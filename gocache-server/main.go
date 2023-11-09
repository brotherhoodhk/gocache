package main

import (
	"encoding/xml"
	"fmt"
	"gocache/body"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"

	"github.com/oswaldoooo/octools/toolsbox"
)

var listenport = 8001
var ROOTPATH = "."
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
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		listener.Close()
		fmt.Println("\ngocache exit with safety")
		os.Exit(1)
	}()
	for {
		con, err := listener.Accept()
		if err != nil {
			errorlog.Println(err)
		} else {
			go body.Process_V2(con)
		}
		// go body.Process(con)

	}
}
