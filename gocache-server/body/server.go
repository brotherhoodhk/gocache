package body

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

var wrbuffsize = 1 * MB
var cachebuffsize = 1 * GB

func Process(Con net.Conn) {
	defer fmt.Println("connection not connected,process exit!")
	var msg = &Message{}
	var rsp = &ReplayStatus{}
	var buff = make([]byte, wrbuffsize)
	var err error
	var langth int
	read := bufio.NewReader(Con)
	for {
		langth, err = read.Read(buff)
		if err != nil {
			errorlog.Println(err.Error())
			Con.Close()
			return
		} else {
			err = json.Unmarshal(buff[:langth], msg)
			if err != nil {
				//the data broken,need client resend data
				// errorlog.Println(err)
			} else {
				buffbytes, code, err := processmsg(msg)
				if code == 200 {
					rsp.Content = buffbytes
				} else {
					rsp.Content = []byte(err.Error())
				}
				rsp.StatusCode = code
				if !sendreply(Con, rsp) {
					errorlog.Println("send reply failed")
					return
				}
			}
		}
	}
}
func Process_V2(con io.ReadWriteCloser) {
	var (
		code uint8
		p    []byte
		err  error
	)
	for {
		code, p, err = Read_V2(con)
		if err == nil {
			code, p, err = processmsg_v2(code, p)
			if err == nil {
				err = Write_V2(con, code, p)
			}
		}
		if err != nil {
			break
		}
	}
	con.Close()
}
func sendreply(Con net.Conn, resp *ReplayStatus) bool {
	resbytes, err := json.Marshal(resp)
	if err != nil {
		errorlog.Println(err)
		fmt.Println("marshal response status failed", err.Error())
		Con.Close()
		return false
	}
	_, err = Con.Write(resbytes)
	if err != nil {
		errorlog.Println(err)
		fmt.Println("write to client failed", err.Error())
		Con.Close()
		return false
	}
	return true
}
