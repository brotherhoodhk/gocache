package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	driver_tools_v2 "github.com/oswaldoooo/gocache-driver/v2"
	"github.com/oswaldoooo/octools/toolsbox"
)

const (
	DEFAULT_DB = "origin_data"
)

var address string
var ROOTPATH = "."
var default_db = DEFAULT_DB

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
	// acp := new(ReplayStatus)
	acp_v2 := make(map[string]any)
	mes.DB = default_db
	// var buff = make([]byte, 10*KB)
	for {
		mes.Act = -1
		mes.Key = ""
		mes.Value = nil
		// acp = &ReplayStatus{}
		fmt.Print("console-> ")
		read := bufio.NewReader(os.Stdout)
		msg, _ := read.ReadString('\n')
		msg = strings.TrimSpace(msg)
		//以下命令不发送,为本地命令
		msgarr := strings.Split(msg, " ")
		switch msgarr[0] {
		case "use":
			if len(msgarr) == 2 {
				mes.DB = msgarr[1]
			}
			goto passthroug
		case "show":
			if len(msgarr) == 2 && msgarr[1] == "db" {
				fmt.Println("current database", mes.DB)
			}
			goto passthroug
		case "clear":
			if len(msgarr) == 1 {
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				err = cmd.Run()

				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("unknown command")
			}
			goto passthroug
		}
		err = Pack(msg, mes)
		if err != nil {
			fmt.Println(err)
		} else {
			acp_v2, err = Do(con, mes)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				for k, ele := range acp_v2 {
					fmt.Println(k, ele)
				}
			}
			// sendbytes, err := json.Marshal(mes)
			// if err != nil {
			// 	fmt.Println(err)
			// } else {
			// senddata:
			// 	_, err = con.Write(sendbytes)
			// 	if err != nil {
			// 		fmt.Println("\n[error] lost connection with host")
			// 		os.Exit(1)
			// 	}
			// 	ReadReply(con, buff, acp)
			// 	switch acp.StatusCode {
			// 	case 500:
			// 		goto senddata
			// 	case 200:
			// 		if mes.Act == 2 || mes.Act == 21 {
			// 			fmt.Println(string(acp.Content))
			// 		}
			// 		continue
			// 	case 400:
			// 		fmt.Println("error 400", string(acp.Content))
			// 	}
			// }

			//v2 edition
			err = driver_tools_v2.WriteTo(con, uint8(mes.Act), mes)
			if err == nil {
				err = driver_tools_v2.ReadFrom(con, &acp_v2)
				if err == nil {

				} else if err == io.EOF {
					//close connection
				} else {
					fmt.Fprintln(os.Stderr, "\n[error]", err.Error())
				}
			}
		}
	passthroug:
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
	case "match":
		if len(msgarr) != 2 {
			return fmt.Errorf("args not correct")
		}
		mes.Key = msgarr[1]
		mes.Act = 21
	case "delete":
		if len(msgarr) < 2 {
			return fmt.Errorf("args not correct")
		}
		mes.Key = strings.Join(msgarr[1:], " ")
		mes.Act = 3
	case "create":
		if len(msgarr) == 2 {
			mes.DB = msgarr[1]
			mes.Act = 10
		} else {
			return fmt.Errorf("args not correct")
		}
	case "save":
		if len(msgarr) == 1 {
			mes.Act = 20
		} else {
			return fmt.Errorf("args not correct")
		}
	case "drop":
		if len(msgarr) == 1 && mes.DB != DEFAULT_DB {
			mes.Act = 30
		} else if len(msgarr) == 1 && mes.DB == DEFAULT_DB {
			return fmt.Errorf("cant delete main db")
		} else {
			return fmt.Errorf("unknown command")
		}
	case "exit":
		fmt.Println("bye")
		os.Exit(1)
	default:
		fmt.Println("unknown command")
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
		rsp.StatusCode = 400
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
	DB    string `json:"db"`
	Key   string `json:"key"`
	Value []byte `json:"value"`
	Act   int    `json:"act"`
}

// status replay
type ReplayStatus struct {
	Content    []byte `json:"content"`
	StatusCode int    `json:"code"`
	Type       string `json:"type"`
}

// services
func Do(in io.ReadWriter, mes *Message) (map[string]any, error) {
	var (
		ans map[string]any = make(map[string]any)
		err error
	)
	err = driver_tools_v2.WriteTo(in, uint8(mes.Act), mes)
	if err == nil {
		err = driver_tools_v2.ReadFrom(in, &ans)
	}
	return ans, err
}
