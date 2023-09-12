package main

import (
	"encoding/xml"
	"fmt"
	"github.com/oswaldoooo/octools/toolsbox"
	"gocache/body"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

var listenport = 8001
var ROOTPATH = os.Getenv("GOCACHE_HOME")
var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")

type confinfo struct {
	XMLName xml.Name `xml:"gocache"`
	Port    int      `xml:"port"`
}

func init() {
	fmt.Println("==========start init external interface==========")
	content, err := ioutil.ReadFile(ROOTPATH + "/conf/conf.xml")
	if err == nil {
		cf := new(confinfo)
		err = xml.Unmarshal(content, cf)
		if err == nil {
			if cf.Port > 0 {
				listenport = cf.Port
			} else {
				fmt.Println("conf port is not correct")
			}
		} else {
			fmt.Println("read conf.xml failed", err)
		}
	} else {
		fmt.Println("open conf.xml failed", err)
	}
}
func main() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(listenport))
	if err != nil {
		fmt.Println(err)
		return
	}
	go http.ListenAndServe(":9999", nil)
	fmt.Println("start listen at ", listenport)
	for {
		con, err := listener.Accept()
		if err != nil {
			errorlog.Println(err)
		}
		defer con.Close()
		go body.Process(con)
	}
}
