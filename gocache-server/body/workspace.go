package body

import (
	"fmt"
	"gocache/basic"
	"strconv"
	"strings"
)

func SaveValue(key string, value string, typeinfo string) error {
	if len(key) > keymaxlength {
		key = key[:keymaxlength]
	}
	switch typeinfo {
	case "string":
		cellmap[key] = basic.NewCell(&String{value})
	case "integer":
		num, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf(value, "is not integer number")
		}
		cellmap[key] = basic.NewCell(&Integer{num})
	case "float":
		num, err := strconv.ParseFloat(value, 10)
		if err != nil {
			return fmt.Errorf(value, "is not float number")
		}
		cellmap[key] = basic.NewCell(&Float{num})
	case "bool":
		value = strings.ToLower(value)
		if value == "true" || value == "false" {
			if value == "true" {
				cellmap[key] = basic.NewCell(&Boolean{true})
			} else {
				cellmap[key] = basic.NewCell(&Boolean{false})
			}
		} else {
			return fmt.Errorf(value, "is not boolean")
		}
	}
	if len(cellmap) > BaseMapSize {
		BaseMapSize = 2 * BaseMapSize
		Save(cellmap)
	}
	return nil
}
func SetKeyValue(key string, value string) {
	typeinfo := findtype(value)
	err := SaveValue(key, value, typeinfo)
	if err != nil {
		fmt.Println(err)
	}
}
func GetKey(key string) []byte {
	if cell, ok := cellmap[key]; ok {
		return cell.GetValue()
	} else {
		return nil
	}
}
func GetAllKeys() []byte {
	res := []byte{}
	for k, _ := range cellmap {
		newres := fmt.Sprintf("%-60v %v", k, string(GetKey(k)))
		res = append(res, []byte(newres)...)
		res = append(res, '\n')
	}
	return res
}
func DeleteKey(key string) {
	if _, ok := cellmap[key]; ok {
		delete(cellmap, key)
	}
	if BaseMapSize > 5 && len(cellmap) < BaseMapSize/2 {
		if BaseMapSize/2 > 5 {
			BaseMapSize = BaseMapSize / 2
		} else {
			BaseMapSize = 5
		}
		Save(cellmap)
	}
}
func ClearAllKeys() {
	cellmap = make(map[string]*basic.Cell)
	Save(cellmap)
	BaseMapSize = 5
}
