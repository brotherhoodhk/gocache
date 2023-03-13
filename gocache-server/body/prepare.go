package body

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/oswaldoooo/octools/toolsbox"
	"io/ioutil"
	"os"
	"strings"
)

var ballast = make([]byte, 300*MB)
var globaldb = NewCustomDB("origin_data")
var ROOTPATH = os.Getenv("GOCACHE_HOME")
var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")

var datapath = ROOTPATH + "/data/"
var keymaxlength = 60
var customdb = make(map[string]*CustomDb) //db
func init() {
	fmt.Println("==========start init main zone==========")
	_, err := os.Stat(ROOTPATH + "/data")
	if err != nil {
		fmt.Println("data directory dont find")
		os.Exit(1)
	}
	buff := make([]byte, cachebuffsize)
	fe, err := os.OpenFile(ROOTPATH+"/data/origin_data.gc", os.O_RDONLY, 0700)
	if err == nil {
		defer fe.Close()
		lang, err := fe.Read(buff)
		if err != nil {
			errorlog.Println(err)
			fmt.Println("read data from disk failed")
		}
		cellmap := globaldb.Cellmap
		err = json.Unmarshal(buff[:lang], &cellmap)
		if err != nil {
			fmt.Println("main data broken")
			errorlog.Println(err)
		}
		globaldb.Cellmap = cellmap
	}
	if len(globaldb.Cellmap) > 5 {
		globaldb.MapContaierSize = len(globaldb.Cellmap)
	}
	//查看有无用户自定义数据库
	fearr, err := ioutil.ReadDir(ROOTPATH + "/data")
	if err == nil && len(fearr) > 1 {
		for _, v := range fearr {
			if v.Name() != "origin_data.gc" && v.Size() > 2 {
				fe, err = os.OpenFile(ROOTPATH+"/data/"+v.Name(), os.O_RDONLY, 0700)
				if err == nil {
					buff = make([]byte, wrbuffsize)
					read := bufio.NewReader(fe)
					lang, err := read.Read(buff)
					if err == nil {
						realdbname := strings.Replace(v.Name(), ".gc", "", 1)
						newdb := NewCustomDB(realdbname)
						err = json.Unmarshal(buff[:lang], &newdb.Cellmap)
						if err == nil {
							if len(newdb.Cellmap) > newdb.MapContaierSize {
								newdb.MapContaierSize = len(newdb.Cellmap)
							}
							customdb[realdbname] = newdb
						}
					}
				}
			}
		}
	}
	if err != nil {
		errorlog.Println(err)
	}
}
