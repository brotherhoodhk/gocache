package body

import (
	"bufio"
	"encoding/json"
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
	var msg = &Message{}
	var rsp = &ReplayStatus{}
	var buff = make([]byte, wrbuffsize)
	var err error
	var langth int
	for {
		read := bufio.NewReader(Con)
		langth, err = read.Read(buff)
		if err != nil {
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
				go sendreply(Con, rsp)
			}
		}
	}
}
func sendreply(Con net.Conn, resp *ReplayStatus) {
	resbytes, err := json.Marshal(resp)
	if err != nil {
		errorlog.Println(err)
		Con.Close()
		return
	}
	_, err = Con.Write(resbytes)
	if err != nil {
		errorlog.Println(err)
		Con.Close()
	}
	return
}
