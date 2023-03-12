package body

import (
	"encoding/json"
	"fmt"
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
func Save(dbinfo *CustomDb) {
	res, err := json.Marshal(&dbinfo.Cellmap)
	if err != nil {
		errorlog.Println(err)
		return
	}
	fe, err := os.OpenFile(datapath+dbinfo.Name+".gc", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
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
func getDB(dbname string) (*CustomDb, error) {
	if len(dbname) < 1 {
		return nil, fmt.Errorf("database name is empty")
	} else if _, ok := customdb[dbname]; !ok && dbname != "origin_data" {
		return nil, fmt.Errorf(dbname, "dont exist")
	} else if dbname == "origin_data" {
		return globaldb, nil
	} else if ve, ok := customdb[dbname]; ok {
		return ve, nil
	} else {
		return nil, fmt.Errorf("unknown error")
	}
}

// 检查数据库在硬盘上是否存在
func checkorigindb(allpath string) bool {
	if _, err := os.Stat(allpath); err == nil {
		return true
	}
	return false
}
func RemoveDBfromDisk(dbinfo *CustomDb) error {
	allpath := datapath + dbinfo.Name + ".gc"
	if checkorigindb(allpath) {
		err := os.Remove(allpath)
		return err
	} else {
		return nil
	}
}

func GetKeyContain(subkey string, dbinfo *CustomDb) (res []byte, err error) {
	resmap := make(map[string][]byte)
	for key, value := range dbinfo.Cellmap {
		if strings.Contains(key, subkey) {
			resmap[key] = value.GetValue()
		}
	}
	if len(resmap) > 0 {
		res, err = json.Marshal(&resmap)
	} else {
		err = fmt.Errorf("no key contain %v", subkey)
	}
	return
}
