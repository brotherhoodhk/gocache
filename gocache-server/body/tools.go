package body

import (
	"encoding/json"
	"fmt"
	"gocache/basic"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
)

// find the value's type
func findtype(value string) (tp string) {
	return basic.Default_Get_Type_Func(value)
}

// save data from cache to disk
func Save(dbinfo *CustomDb) {
	res, err := json.Marshal(&dbinfo.Cellmap)
	if err != nil {
		errorlog.Println(err)
		return
	}
	fe, err := os.OpenFile(datapath+dbinfo.Name+".gc", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		errorlog.Println(err)
		return
	}
	_, err = fe.Write(res)
	if err != nil {
		errorlog.Println(err)
		return
	}
}
func getDB(dbname string) (*CustomDb, error) {
	if len(dbname) < 1 {
		return nil, fmt.Errorf("database name is empty")
	} else if _, ok := customdb[dbname]; !ok && dbname != "origin_data" {
		return nil, fmt.Errorf(dbname, "dont exist")
	} else if dbname == "origin_data" {
		return globaldb, nil
	} else if ve, ok := customdb[dbname]; ok {
		return ve, nil
	} else {
		return nil, fmt.Errorf("unknown error")
	}
}

// 检查数据库在硬盘上是否存在
func checkorigindb(allpath string) bool {
	if _, err := os.Stat(allpath); err == nil {
		return true
	}
	return false
}
func RemoveDBfromDisk(dbinfo *CustomDb) error {
	allpath := datapath + dbinfo.Name + ".gc"
	if checkorigindb(allpath) {
		err := os.Remove(allpath)
		return err
	} else {
		return nil
	}
}

func GetKeyContain(subkey string, dbinfo *CustomDb) (res []byte, err error) {
	resmap := make(map[string][]byte)
	for key, value := range dbinfo.Cellmap {
		if strings.Contains(key, subkey) {
			resmap[key] = value.GetValue()
		}
	}
	if len(resmap) > 0 {
		res, err = json.Marshal(&resmap)
	} else {
		err = fmt.Errorf("no key contain %v", subkey)
	}
	return
}
func GetKeyContain_V2(subkey string, dbinfo *CustomDb) (res map[string][]byte, err error) {
	res = make(map[string][]byte)
	for key, value := range dbinfo.Cellmap {
		if strings.Contains(key, subkey) {
			res[key] = value.GetValue()
		}
	}
	if len(res) == 0 {
		err = fmt.Errorf("no key contain %v", subkey)
	}
	return
}

// v2 edition
func Read_V2(in io.Reader) (code uint8, p []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = io.EOF
		}
	}()
	reader := NewReader(10 << 10)
	code = uint8(reader.Read(in, 1, Bigendian))
	data_len := reader.Read(in, 2, Bigendian)
	if data_len > 0 {
		p = reader.RawRead(in, int(data_len))
	}
	return
}
func Write_V2(out io.Writer, code uint8, content []byte) error {
	var err error
	newcontent := make([]byte, 3+len(content))
	newcontent[0] = code
	err = Bigendian.Write(uint64(len(content)), 2, newcontent[1:3])
	if err == nil {
		if len(content) > 0 {
			copy(newcontent[3:], content)
		}
		_, err = out.Write(newcontent)
	}
	return err
}

// utils file
var (
	Bigendian    = new(BigEndian)
	Littleendian = new(LittleEndian)
)

type endian interface {
	Uint(src []byte, n uint) uint64
	Write(v uint64, n int, p []byte) error
}

type BigEndian struct {
}
type LittleEndian struct {
}

func (s *BigEndian) Uint(src []byte, n uint) uint64 {
	step := n
	var length uint = uint(len(src))
	if length < step {
		step = length
	}
	var ans uint64
	if step > 0 {
		var i uint
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintln(os.Stderr, "error pos", i, "src length", len(src))
				os.Exit(1)
			}
		}()
		for i = 0; i < step; i++ {
			ans *= 256
			ans += uint64(src[i])
		}
	}
	return ans
}
func (s *BigEndian) Write(v uint64, n int, p []byte) error {
	if len(p) < n {
		return Str_Error("p out of range n")
	}
	var i int
	for i < n {
		p[n-1-i] = uint8(v % 256)
		if v > 0 {
			v /= 256
		}
		i++
	}
	return nil
}

func (s *LittleEndian) Uint(src []byte, n uint) uint64 {
	step := n
	var length uint = uint(len(src))
	if length < step {
		step = length
	}
	var ans uint64
	if step > 1 {
		var i int
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintln(os.Stderr, "error pos", i, "src length", len(src), "step", step)
				os.Exit(1)
			}
		}()
		for i = int(step - 1); i >= 0; i-- {
			ans *= 256
			ans += uint64(src[i])
		}
	}
	return ans
}
func (s *LittleEndian) Write(v uint64, n int, p []byte) error {
	if len(p) < n {
		return Str_Error("p out of range n")
	}
	var i int
	for i < n {
		p[i] = uint8(v % 256)
		if v > 0 {
			v /= 256
		}
		i++
	}
	return nil
}

type Reader struct {
	cache_buffer []byte
	mux          sync.Mutex
}

func (s *Reader) Read(reader io.Reader, n int, end endian) uint64 {
	s.mux.Lock()
	defer s.mux.Unlock()
	// lang, err := reader.Read(s.cache_buffer[0:n])
	err := read(reader, s.cache_buffer[0:n])
	if err == nil {
		if n == 1 {
			return uint64(s.cache_buffer[0])
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintln(os.Stderr, "[panic error] lang", n, len(s.cache_buffer), r)
				os.Exit(1)
			}
		}()
		return end.Uint(s.cache_buffer[0:n], uint(n))
	} else {
		panic(err)
	}
	return 0
}
func (s *Reader) RawRead(reader io.Reader, n int) []byte {
	s.mux.Lock()
	defer s.mux.Unlock()
	var (
		ans []byte = nil
	)
	err := read(reader, s.cache_buffer[:n])
	if err == nil {
		ans = make([]byte, n)
		copy(ans, s.cache_buffer[:n])
	} else {
		panic(err)
	}
	return ans
}
func NewReader(size uint64) *Reader {
	return &Reader{cache_buffer: make([]byte, size)}
}

type Str_Error string

func (s Str_Error) Error() string {
	return string(s)
}

type ReadWriter int

func (s ReadWriter) Read(p []byte) (n int, err error) {
	return syscall.Read(int(s), p)
}

func (s ReadWriter) Write(p []byte) (n int, err error) {
	return syscall.Write(int(s), p)
}

// read full byte array
func read(in io.Reader, p []byte) error {
	var (
		err   error
		n     int
		start int
		lang  = len(p)
	)
	n, err = in.Read(p[start:])
	for err == nil && start < lang {
		start += n
		if start >= lang {
			break
		}
		n, err = in.Read(p[start:])
		if n == 0 {
			err = io.EOF
		}
	}
	return err
}
