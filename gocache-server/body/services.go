package body

import (
	"fmt"
	"gocache/basic"

	"github.com/oswaldoooo/octools/toolsbox"
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
