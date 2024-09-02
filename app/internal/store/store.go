package store

import (
	"fmt"
	"time"
)

type data struct {
	value string;
	createdAt uint;
	expireAt uint;
}
var store = make(map[string]data)

func SetStore(key, val string, args ...uint) {
	timestamp := uint(time.Now().UnixMilli())

	if len(args) > 0 {
		store[key] = data {
			value: val,
			createdAt: timestamp,
			expireAt: timestamp + args[0],
		}
	} else {
		store[key] = data {
			value: val,
			createdAt: timestamp,
		}
	}
}

func GetStore(key string) (string, bool) {
	val, isPresent := store[key]
	fmt.Println(val.expireAt)
	fmt.Println(uint(time.Now().UnixMilli()))
	if val.expireAt != 0 {
		if uint(time.Now().UnixMilli()) >= val.expireAt {
			delete(store, key)
			return "", false
		}
	}

	return val.value, isPresent
}