package body

import (
	"encoding/json"
	"fmt"
	"gocache/basic"
	"gocache/utils"

	"github.com/oswaldoooo/octools/toolsbox"
)

const ( //status code
	OK    = 0x01
	ERROR = 0x02
)

const ( //commuicate signal
	PING = 0x10
)

var (
	Encode func(any) ([]byte, error) = json.Marshal
	Decode func([]byte, any) error   = json.Unmarshal
)
var processlog = toolsbox.LogInit("process", ROOTPATH+"/logs/process.log")

func processmsg(msg *Message) ([]byte, int, error) {
	switch msg.Act {
	case 1, 31:
		var dbinfo *CustomDb
		if len(msg.Value) < 1 {
			return nil, 400, fmt.Errorf("value is empty")
		} else if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		SetKeyValue(msg.Key, string(msg.Value), dbinfo)
		return nil, 200, nil
	case 11: //compare and set key,return version
		var dbinfo *CustomDb
		if len(msg.Value) < 1 {
			return nil, 400, fmt.Errorf("value is empty")
		} else if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		var version_id uint32 = 0
		content := GetKey(msg.Key, dbinfo)
		if content == nil {
			content = make([]byte, len(msg.Value)+2)
			content[0] = 0
			content[1] = 0
			copy(content[2:], []byte(msg.Value))
			SetKeyValue(msg.Key, string(content), dbinfo)
		} else if len(content) >= 2 && string(content[2:]) != string(msg.Value) {
			//update version
			version_id = uint32(content[0])*256 + uint32(content[1])
			version_id++
			content[0] = byte(version_id / 256)
			content[1] = byte(version_id % 256)
			copy(content[2:], msg.Value)
			SetKeyValue(msg.Key, string(content), dbinfo)
		} else {
			return nil, 400, fmt.Errorf("msg old value is not compareable")
		}
		return []byte{byte(version_id / 256), byte(version_id % 256)}, 200, nil
	case 2: //get key
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		if msg.Key == "*" {
			return GetAllKeys(dbinfo), 200, nil
		}
		return GetKey(msg.Key, dbinfo), 200, nil
	case 32:
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		}
		if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		// fmt.Println(string(GetKeys(msg.Key, dbinfo)))
		return GetKeys(msg.Key, dbinfo), 200, nil
	case 322:
		//驱动接口获得全部键
		var dbinfo *CustomDb
		if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		return GetAllKeysInterface(dbinfo), 200, nil
	case 21:
		var dbinfo *CustomDb
		//模糊查询 fuzzy query
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		return FuzzyMatch(msg.Key, dbinfo), 200, nil
	case 341:
		//find the key that contain rune
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		resbytes, err := GetKeyContain(msg.Key, dbinfo)
		if err != nil {
			return nil, 400, err
		}
		return resbytes, 200, nil
	case 3, 33:
		//delete keys
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			return nil, 400, err
		}
		if msg.Key == "*" {
			ClearAllKeys(dbinfo)
		} else {
			DeleteKeys(msg.Key, dbinfo)
		}
		return nil, 200, nil
	case 10, 310:
		//create db
		var err error
		if len(msg.DB) > 0 {
			if CreateDB(msg.DB) {
				//debug line
				// processlog.Println("update db,", customdb)
				return nil, 200, nil
			} else {
				err = fmt.Errorf(msg.DB, "already exist")
			}
		} else {
			err = fmt.Errorf("dbname is empty")
		}
		return nil, 400, err
	case 20, 320:
		//手动储存
		var err error
		if dbinfo, err := getDB(msg.DB); err == nil {
			if len(dbinfo.Cellmap) > 0 {
				Save(dbinfo)
				return nil, 200, nil
			} else {
				err = fmt.Errorf("database is empty,please put some data into it first")
			}
		}
		return nil, 400, err
	case 30, 330:
		//删除数据库
		var err error
		if dbinfo, err := getDB(msg.DB); err == nil {
			delete(customdb, dbinfo.Name)
			err = RemoveDBfromDisk(dbinfo)
			if err == nil {
				return nil, 200, nil
			}
		}
		return nil, 400, err
	case 90:
		//加载配置文件
		var err error
		if len(msg.Key) < 5 {
			err = fmt.Errorf("site configure file name is not correct")
		} else {
			err = LoadConf(msg.Key)
		}
		if err != nil {
			return nil, 400, err
		}
		return nil, 200, nil
	case 901:
		//load plugin list
		basic.LoadPluginList()
		return nil, 200, nil
	default:
		if basic.IsOpenExtensions {
			if addfunc, ok := basic.AddtionalCommand[msg.Act]; ok {
				return addfunc(msg.Key, msg.DB, msg.Value)
			}
		}
	}
	return nil, 400, fmt.Errorf("unknown command")
}

/*
v2 signal changelog: old=new
322 =34
341=35
310=36
320=37
330=38
910=39
*/
func processmsg_v2(code uint8, rawcontent []byte) (ans_code uint8, p []byte, err error) {
	var (
		ans_map map[string]any = make(map[string]any)
		msg     Message
	)
	if len(rawcontent) > 0 {
		err = Decode(rawcontent, &msg)
		if err != nil {
			errorlog.Println("decode failed", err.Error())
			return
		}
	}
	switch code {
	case PING:
		if len(rawcontent) > 0 { //it's pingx
			err = pingx(rawcontent)
			if err == nil {
				ans_code = OK
			} else {
				ans_code = ERROR
			}
		}
	case 1, 31:
		var dbinfo *CustomDb
		if len(msg.Value) < 1 {
			ans_code = ERROR
			ans_map["error"] = "value is empty"
		} else if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
		SetKeyValue(msg.Key, string(msg.Value), dbinfo)
		ans_code = OK
	case 11: //compare and set key,return version
		var dbinfo *CustomDb
		if len(msg.Value) < 1 {
			ans_code = ERROR
			ans_map["error"] = "value is empty"
		} else if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
		var version_id uint32 = 0
		content := GetKey(msg.Key, dbinfo)
		if content == nil {
			content = make([]byte, len(msg.Value)+2)
			content[0] = 0
			content[1] = 0
			copy(content[2:], []byte(msg.Value))
			SetKeyValue(msg.Key, string(content), dbinfo)
		} else if len(content) >= 2 && string(content[2:]) != string(msg.Value) {
			//update version
			version_id = uint32(content[0])*256 + uint32(content[1])
			version_id++
			content[0] = byte(version_id / 256)
			content[1] = byte(version_id % 256)
			copy(content[2:], msg.Value)
			SetKeyValue(msg.Key, string(content), dbinfo)
			ans_code = OK
			ans_map["version"] = version_id
		} else {
			ans_code = ERROR
			ans_map["error"] = "msg old value is not compareable"
		}
	case 2: //get key
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
		if msg.Key == "*" {
			ans_map = utils.TransMap(GetAllKeys_V2(dbinfo))
		} else {
			ans_map[msg.Key] = GetKey(msg.Key, dbinfo)
		}
		ans_code = OK
	case 32: //not implment
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			ans_code = OK
			dbinfo = dbinfocopy
			ans_map = utils.TransMap(GetKeys_V2(msg.Key, dbinfo))
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
		errorlog.Println(err, len(ans_map))
	case 34:
		//驱动接口获得全部键
		var dbinfo *CustomDb
		if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}

		ans_map = utils.TransMap(GetAllKeysInterface_V2(dbinfo))
		ans_code = OK
	case 21:
		var dbinfo *CustomDb
		//模糊查询 fuzzy query
		if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
		ans_code = OK
		ans_map = utils.TransMap(FuzzyMatch_V2(msg.Key, dbinfo))
	case 35:
		//find the key that contain rune
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
			resmap, err := GetKeyContain_V2(msg.Key, dbinfo)
			if err != nil {
				ans_code = ERROR
				ans_map["error"] = err.Error()
			} else {
				ans_code = OK
				ans_map = utils.TransMap(resmap)
			}
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
	case 3, 33:
		//delete keys
		var dbinfo *CustomDb
		if len(msg.Key) < 1 {
			ans_code = ERROR
			ans_map["error"] = "key is empty"
		} else if dbinfocopy, err := getDB(msg.DB); err == nil {
			dbinfo = dbinfocopy
			if msg.Key == "*" {
				ClearAllKeys(dbinfo)
			} else {
				DeleteKeys(msg.Key, dbinfo)
			}
		} else {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
		ans_code = OK
	case 10, 36:
		//create db
		var err error
		if len(msg.DB) > 0 {
			if CreateDB(msg.DB) {
				//debug line
				// processlog.Println("update db,", customdb)
				ans_code = OK
			} else {
				err = fmt.Errorf(msg.DB, "already exist")
			}
		} else {
			err = fmt.Errorf("dbname is empty")
		}
		if err != nil {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}

	case 20, 37:
		//手动储存
		var err error
		if dbinfo, err := getDB(msg.DB); err == nil {
			if len(dbinfo.Cellmap) > 0 {
				Save(dbinfo)
				ans_code = OK
			} else {
				err = fmt.Errorf("database is empty,please put some data into it first")
			}
		}
		if err != nil {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}

	case 30, 38:
		//删除数据库
		var err error
		if dbinfo, err := getDB(msg.DB); err == nil {
			delete(customdb, dbinfo.Name)
			err = RemoveDBfromDisk(dbinfo)
			if err == nil {
				ans_code = OK
			}
		}
		if err != nil {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		}
	case 90:
		//加载配置文件
		var err error
		if len(msg.Key) < 5 {
			err = fmt.Errorf("site configure file name is not correct")
		} else {
			err = LoadConf(msg.Key)
		}
		if err != nil {
			ans_code = ERROR
			ans_map["error"] = err.Error()
		} else {
			ans_code = OK
		}

	case 39:
		//load plugin list
		basic.LoadPluginList()
		ans_code = OK
	default:
		if basic.IsOpenExtensions {
			if addfunc, ok := basic.AddtionalCommand_V2[code]; ok {
				ans_map, err = addfunc(msg.Key, msg.DB, msg.Value)
				if err != nil {
					ans_code = ERROR
					ans_map["error"] = err.Error()
				} else {
					ans_code = OK
				}
			} else {
				ans_code = ERROR
				ans_map["error"] = "no such command"
			}
		} else {
			ans_code = ERROR
			ans_map["error"] = "no such command"
		}
	}
	if len(ans_map) > 0 {
		p, err = Encode(&ans_map)
		if err != nil {
			errorlog.Println("[encode failed]", err.Error())
		}
	}
	return
}
func pingx(p []byte) error {
	var (
		err  error
		data map[string]string = make(map[string]string)
	)
	err = Decode(p, &data)
	if err == nil {
		//verify whether can visit database
		passwd := data["passwd"]
		db := data["database"]
		if utils.VerifyPermit(passwd, db) {
			return nil
		} else {
			return Str_Error("permission denied")
		}
	} else {
		errorlog.Println("parsed error", err.Error())
	}
	return err
}
