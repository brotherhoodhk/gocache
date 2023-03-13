package main

import (
	"encoding/xml"
	"fmt"
	"gocache/body"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/oswaldoooo/octools/toolsbox"
)

var listenport = 8001
var ROOTPATH = os.Getenv("GOCACHE_HOME")
var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")

type confinfo struct {
	XMLName xml.Name `xml:"gocache"`
	Port    int      `xml:"port"`
}

func init() {
	// baselist, err := toolsbox.ParseList(ROOTPATH + "/conf/site.cnf")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
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
	// if port, ok := baselist["port"]; ok {
	// 	realport, err := strconv.Atoi(port)
	// 	if err == nil {
	// 		listenport = realport
	// 	} else {
	// 		fmt.Println(port, "is not correct number,gocache will use default port")
	// 	}
	// } else {
	// 	fmt.Println("site not configure listen port,gocache will use default port")
	// }
}
func main() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(listenport))
	if err != nil {
		fmt.Println(err)
		return
	}
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
