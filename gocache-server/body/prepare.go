package body

import (
	"encoding/json"
	"fmt"
	"gocache/basic"
	"os"

	"github.com/oswaldoooo/octools/toolsbox"
)

var ROOTPATH = os.Getenv("GOCACHE_HOME")
var cellmap = make(map[string]*basic.Cell)
var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")
var BaseMapSize = 5
var datapath = ROOTPATH + "/data/"
var keymaxlength = 60

func init() {
	_, err := os.Stat(ROOTPATH + "/data")
	if err != nil {
		// fmt.Println("data directory dont find")
		os.Exit(1)
	}
	buff := make([]byte, cachebuffsize)
	fe, err := os.OpenFile(ROOTPATH+"/data/origin_data.gc", os.O_RDONLY, 0700)
	if err == nil {
		lang, err := fe.Read(buff)
		if err != nil {
			errorlog.Println(err)
			fmt.Println("read data from disk failed")
		}
		err = json.Unmarshal(buff[:lang], &cellmap)
		if err != nil {
			fmt.Println("data broken")
			errorlog.Println(err)
		}
	}
	if len(cellmap) > 5 {
		BaseMapSize = len(cellmap)
	}
}
