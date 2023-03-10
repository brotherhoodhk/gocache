package body

import (
	"fmt"

	"github.com/oswaldoooo/octools/toolsbox"
)

var processlog = toolsbox.LogInit("process", ROOTPATH+"/logs/process.log")

func processmsg(msg *Message) ([]byte, int, error) {
	switch msg.Act {
	case 1:
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
	case 2:
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
	case 3:
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
	case 10:
		//create db
		var err error
		if len(msg.DB) > 0 {
			if CreateDB(msg.DB) {
				//debug line
				processlog.Println("update db,", customdb)
				return nil, 200, nil
			} else {
				err = fmt.Errorf(msg.DB, "already exist")
			}
		} else {
			err = fmt.Errorf("dbname is empty")
		}
		return nil, 400, err
	case 20:
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
	case 30:
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
	}
	return nil, 400, nil
}
