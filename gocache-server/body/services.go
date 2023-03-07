package body

import "fmt"

func processmsg(msg *Message) ([]byte, int, error) {
	switch msg.Act {
	case 1:
		if len(msg.Value) < 1 {
			return nil, 400, fmt.Errorf("value is empty")
		} else if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		}
		SetKeyValue(msg.Key, string(msg.Value))
		return nil, 200, nil
	case 2:
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		}
		if msg.Key == "*" {
			return GetAllKeys(), 200, nil
		}
		return GetKey(msg.Key), 200, nil
	case 3:
		if len(msg.Key) < 1 {
			return nil, 400, fmt.Errorf("key is empty")
		}
		if msg.Key == "*" {
			ClearAllKeys()
		} else {
			DeleteKey(msg.Key)
		}
		return nil, 200, nil
	}
	return nil, 400, nil
}
