package body

import (
	"encoding/json"
	"fmt"
	"gocache/basic"
	"strconv"
	"strings"
)

func SaveValue(key string, value string, typeinfo string, dbinfo *CustomDb) error {
	if len(key) > keymaxlength {
		key = key[:keymaxlength]
	}
	cellmap := &dbinfo.Cellmap
	BaseMapSize := &dbinfo.MapContaierSize
	dbinfo.Mutex.Lock()
	defer dbinfo.Mutex.Unlock()
	switch typeinfo {
	case "string", "str":
		(*cellmap)[key] = basic.NewCell(&String{value})
	case "integer", "int":
		num, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf(value, "is not integer number")
		}
		(*cellmap)[key] = basic.NewCell(&Integer{num})
	case "float":
		num, err := strconv.ParseFloat(value, 10)
		if err != nil {
			return fmt.Errorf(value, "is not float number")
		}
		(*cellmap)[key] = basic.NewCell(&Float{num})
	case "bool", "boolean":
		value = strings.ToLower(value)
		if value == "true" || value == "false" {
			if value == "true" {
				(*cellmap)[key] = basic.NewCell(&Boolean{true})
			} else {
				(*cellmap)[key] = basic.NewCell(&Boolean{false})
			}
		} else {
			return fmt.Errorf(value, "is not boolean")
		}
	}
	// dbinfo.Cellmap = cellmap
	if len(*cellmap) > *BaseMapSize {
		(*BaseMapSize) = 2 * (*BaseMapSize)
		Save(dbinfo)
	}
	// dbinfo.MapContaierSize = BaseMapSize
	return nil
}
func SetKeyValue(key string, value string, dbinfo *CustomDb) {
	typeinfo := findtype(value)
	err := SaveValue(key, value, typeinfo, dbinfo)
	if err != nil {
		fmt.Println(err)
	}
}
func GetKey(key string, dbinfo *CustomDb) []byte {
	if cell, ok := dbinfo.Cellmap[key]; ok {
		return cell.GetValue()
	} else {
		return nil
	}
}

// return string slice
func GetKeys(keys string, dbinfo *CustomDb) (resbytes []byte) {
	resstr := []string{}
	if strings.ContainsRune(keys, ' ') {
		keysarr := strings.Split(keys, " ")
		var buffres []byte
		dbinfo.Mutex.RLock()
		for _, ve := range keysarr {
			if len(ve) > 0 {
				buffres = GetKey(ve, dbinfo)
				if buffres != nil {
					resstr = append(resstr, string(buffres))
				}
			}
		}
		dbinfo.Mutex.RUnlock()
	} else {
		resstr = append(resstr, string(GetKey(keys, dbinfo)))
	}
	tibytes, err := json.Marshal(&resstr)
	if err == nil {
		resbytes = tibytes
	}
	return
}
func GetKeys_V2(keys string, dbinfo *CustomDb) (resbytes []byte) {
	panic("not implment")
	// resstr := []string{}
	// if strings.ContainsRune(keys, ' ') {
	// 	keysarr := strings.Split(keys, " ")
	// 	var buffres []byte
	// 	dbinfo.Mutex.RLock()
	// 	for _, ve := range keysarr {
	// 		if len(ve) > 0 {
	// 			buffres = GetKey(ve, dbinfo)
	// 			if buffres != nil {
	// 				resstr = append(resstr, string(buffres))
	// 			}
	// 		}
	// 	}
	// 	dbinfo.Mutex.RUnlock()
	// } else {
	// 	resstr = append(resstr, string(GetKey(keys, dbinfo)))
	// }
	// tibytes, err := json.Marshal(&resstr)
	// if err == nil {
	// 	resbytes = tibytes
	// }
	return
}
func GetAllKeysInterface(dbinfo *CustomDb) (res []byte) {
	resmap := make(map[string][]byte)
	dbinfo.Mutex.RLock()
	for k := range dbinfo.Cellmap {
		resmap[k] = GetKey(k, dbinfo)
	}
	dbinfo.Mutex.RUnlock()
	res, _ = json.Marshal(&resmap)
	return
}
func GetAllKeysInterface_V2(dbinfo *CustomDb) map[string][]byte {
	resmap := make(map[string][]byte)
	dbinfo.Mutex.RLock()
	for k := range dbinfo.Cellmap {
		resmap[k] = GetKey(k, dbinfo)
	}
	dbinfo.Mutex.RUnlock()
	return resmap
}
func GetAllKeys(dbinfo *CustomDb) []byte {
	res := []byte{}
	dbinfo.Mutex.RLock()
	defer dbinfo.Mutex.RUnlock()
	for k := range dbinfo.Cellmap {
		newres := fmt.Sprintf("%-60v %v", k, string(GetKey(k, dbinfo)))
		res = append(res, []byte(newres)...)
		res = append(res, '\n')
	}
	return res
}
func GetAllKeys_V2(dbinfo *CustomDb) map[string][]byte {
	res := make(map[string][]byte)
	dbinfo.Mutex.RLock()
	defer dbinfo.Mutex.RUnlock()
	for k := range dbinfo.Cellmap {
		res[k] = GetKey(k, dbinfo)
	}
	return res
}

// delete single key
func DeleteKey(key string, dbinfo *CustomDb) {
	cellmap := &dbinfo.Cellmap
	BaseMapSize := dbinfo.MapContaierSize
	dbinfo.Mutex.Lock()
	defer dbinfo.Mutex.Unlock()
	if _, ok := (*cellmap)[key]; ok {
		delete((*cellmap), key)
	}
	// dbinfo.Cellmap = cellmap
	if BaseMapSize > 5 && len(*cellmap) < BaseMapSize/2 {
		if BaseMapSize/2 > 5 {
			BaseMapSize = BaseMapSize / 2
		} else {
			BaseMapSize = 5
		}
		Save(dbinfo)
	}
	dbinfo.MapContaierSize = BaseMapSize
}

// delete multipe keys
func DeleteKeys(keys string, dbinfo *CustomDb) {
	dbinfo.Mutex.Lock()
	defer dbinfo.Mutex.Unlock()
	if strings.ContainsRune(keys, ' ') {
		keysarr := strings.Split(keys, " ")
		for _, v := range keysarr {
			if len(v) > 0 {
				DeleteKey(v, dbinfo)
			}
		}
	} else {
		DeleteKey(keys, dbinfo)
	}
}

// delete all keys
func ClearAllKeys(dbinfo *CustomDb) {
	dbinfo.Mutex.Lock()
	defer dbinfo.Mutex.Unlock()
	dbinfo.Cellmap = make(map[string]*basic.Cell)
	Save(dbinfo)
	dbinfo.MapContaierSize = 5
}
func FuzzyMatch(target string, dbinfo *CustomDb) []byte {
	res := []byte{}
	dbinfo.Mutex.RLock()
	defer dbinfo.Mutex.RUnlock()
	for k := range dbinfo.Cellmap {
		if basic.Default_Fuzzy_Match_Func(target, k) {
			newres := fmt.Sprintf("%-60v %v", k, string(GetKey(k, dbinfo)))
			res = append(res, []byte(newres)...)
			res = append(res, '\n')
		}
	}
	return res
}
func FuzzyMatch_V2(target string, dbinfo *CustomDb) map[string][]byte {
	res := make(map[string][]byte)
	dbinfo.Mutex.RLock()
	defer dbinfo.Mutex.RUnlock()
	for k := range dbinfo.Cellmap {
		if basic.Default_Fuzzy_Match_Func(target, k) {
			res[k] = GetKey(k, dbinfo)
		}
	}
	return res
}

// create new db,if exist,return false
func CreateDB(dbname string) bool {
	if _, ok := customdb[dbname]; !ok {
		customdb[dbname] = NewCustomDB(dbname)
		return true
	} else {
		return false
	}
}

// load the target site configure file
func LoadConf(filename string) (err error) {
	filepath := ROOTPATH + "/conf/" + filename
	err = basic.ReadConf(filepath)
	return
}
