package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/oswaldoooo/octools/toolsbox"
)

var address string
var ROOTPATH = os.Getenv("GOCACHECLI_HOME")

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

func init() {
	basicmap, err := toolsbox.ParseList(ROOTPATH + "/conf/site.cnf")
	if err != nil {
		fmt.Println("read site configure failed")
		os.Exit(-1)
	}
	if realadd, ok := basicmap["gocache_address"]; ok {
		address = realadd
	} else {
		fmt.Println("cant find gocache address in", ROOTPATH+"/conf/site.cnf")
		os.Exit(1)
	}
}
func main() {
	con, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(address, "is offline")
		return
	}
	mes := new(Message)
	acp := new(ReplayStatus)
	var buff = make([]byte, 10*KB)
	for {
		mes = &Message{}
		acp = &ReplayStatus{}
		fmt.Print("console-> ")
		read := bufio.NewReader(os.Stdout)
		msg, _ := read.ReadString('\n')
		msg = strings.TrimSpace(msg)
		err = Pack(msg, mes)
		if err != nil {
			fmt.Println(err)
		} else {
			sendbytes, err := json.Marshal(mes)
			if err != nil {
				fmt.Println(err)
			} else {
			senddata:
				con.Write(sendbytes)
				ReadReply(con, buff, acp)
				switch acp.StatusCode {
				case 500:
					goto senddata
				case 200:
					if mes.Act == 2 {
						fmt.Println(string(acp.Content))
					}
					continue
				case 400:
					fmt.Println("error 400")
				}
			}
		}
	}
}
func checkonline(con net.Conn) {
	for {
		_, err := con.Write(nil)
		if err != nil {
			fmt.Println("connection lost")
			os.Exit(2)
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}
func Pack(msg string, mes *Message) error {
	if len(msg) < 1 {
		return fmt.Errorf("nothing here")
	}
	msgarr := strings.Split(msg, " ")
	if len(msgarr) < 1 {
		return fmt.Errorf("nothing here")
	}
	switch msgarr[0] {
	case "set":
		if len(msgarr) < 3 {
			return fmt.Errorf("args not enough")
		}
		mes.Key = msgarr[1]
		mes.Value = []byte(strings.Join(msgarr[2:], " "))
		mes.Act = 1
	case "get":
		if len(msgarr) != 2 {
			return fmt.Errorf("args not correct")
		}
		mes.Key = msgarr[1]
		mes.Act = 2
	case "delete":
		if len(msgarr) != 2 {
			return fmt.Errorf("args not correct")
		}
		mes.Key = msgarr[1]
		mes.Act = 3
	case "exit":
		fmt.Println("bye")
		os.Exit(1)
	}
	return nil

}
func ReadReply(con net.Conn, buff []byte, rsp *ReplayStatus) {
	lang, err := con.Read(buff)
	if err != nil {
		if err == io.EOF {
			fmt.Println("warn: connection closed")
			os.Exit(2)
		}
		fmt.Println(err)
		con.Close()
		return
	}
	err = json.Unmarshal(buff[:lang], rsp)
	if err != nil {
		rsp.StatusCode = 500
	}

}

// accept msg
type Message struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
	Act   int    `json:"act"`
}

// status replay
type ReplayStatus struct {
	Content    []byte `json:"content"`
	StatusCode int    `json:"code"`
}
