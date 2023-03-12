package main

import (
	"fmt"
	"gocache/body"
	"net"
	"os"
	"strconv"

	"github.com/oswaldoooo/octools/toolsbox"
)

var listenport = 8001
var ROOTPATH = os.Getenv("GOCACHE_HOME")
var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")

func init() {
	baselist, err := toolsbox.ParseList(ROOTPATH + "/conf/site.cnf")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if port, ok := baselist["port"]; ok {
		realport, err := strconv.Atoi(port)
		if err == nil {
			listenport = realport
		} else {
			fmt.Println(port, "is not correct number,gocache will use default port")
		}
	} else {
		fmt.Println("site not configure listen port,gocache will use default port")
	}
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
