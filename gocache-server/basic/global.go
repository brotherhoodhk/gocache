package basic

import (
	"strconv"
	"strings"

	"github.com/oswaldoooo/octools/datastore"
)

var CommonPath string
var IsOpenExtensions = false //是否打开扩展功能池
var CompareRate = 60
var Default_Fuzzy_Match_Func func(string, string) bool = comparesimple
var Default_Get_Type_Func func(string) string = getcontenttype
var AddtionalCommand = make(map[int]func(key string, db string, val []byte) ([]byte, int, error)) //扩展功能池
var Fuzzy_Match_Func_Pool = make(map[string]func(target, tocompare string) bool)
var AddtionalCommand_V2 = make(map[uint8]Method) //扩展功能池
type Method func(key string, db string, val []byte) (map[string]any, error)

func comparesimple(origin, tocompare string) bool {
	return datastore.Comparestr(origin, tocompare, CompareRate)
}
func getcontenttype(value string) string {
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
