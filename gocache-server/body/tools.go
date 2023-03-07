package body

import (
	"encoding/json"
	"gocache/basic"
	"os"
	"strconv"
	"strings"
)

// find the value's type
func findtype(value string) (tp string) {
	var ere error
	if _, ere = strconv.ParseFloat(value, 10); strings.ContainsRune(value, '.') && ere == nil {
		//it's float
		return "float"
	} else if _, ere = strconv.Atoi(value); ere == nil {
		//it's integer
		return "integer"
	} else if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
		//it's boolean
		return "boolean"
	} else {
		return "string"
	}
}

// save data from cache to disk
func Save(data map[string]*basic.Cell) {
	res, err := json.Marshal(data)
	if err != nil {
		errorlog.Println(err)
		return
	}
	fe, err := os.OpenFile(datapath+"origin_data.gc", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		errorlog.Println(err)
		return
	}
	_, err = fe.Write(res)
	if err != nil {
		errorlog.Println(err)
		return
	}
}
