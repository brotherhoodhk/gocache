package main

import (
	"encoding/json"
	"fmt"
	"os"

	driver_tools_v2 "github.com/oswaldoooo/gocache-driver/v2"
)

var db *driver_tools_v2.CacheDB_V2 = driver_tools_v2.NewCacheDB_V2("127.0.0.1", 8001, "", "one")

func init() {
	driver_tools_v2.Encode = json.Marshal
	driver_tools_v2.Decode = json.Unmarshal
}
func main() {
	err := db.Connect()
	if err == nil {
		defer db.Close()
		fmt.Println("connect pass")
		err = db.CreateDB()
		if err == nil {
			fmt.Println("create db pass")
		} else {
			fmt.Fprintln(os.Stderr, "[error]", err.Error())
		}
		err = db.SetKey("halo", "world")
		if err == nil {
			fmt.Println("set key pass")
		} else {
			fmt.Fprintln(os.Stderr, "[error]", err.Error())
		}
		var ansbytes [][]byte
		ansbytes, err = db.GetKeys("halo")
		if err == nil {
			for _, ele := range ansbytes {
				fmt.Println(string(ele))
			}
		} else {
			fmt.Fprintln(os.Stderr, "[error]", err.Error())
		}
		var ans map[string][]byte
		ans, err = db.GetAllKeys()
		if err == nil {
			fmt.Println("get all key pass")
			for key, ele := range ans {
				fmt.Println(key, string(ele))
			}
		} else {
			fmt.Fprintln(os.Stderr, "[error]", err.Error())
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error]", err.Error())
	}
}
