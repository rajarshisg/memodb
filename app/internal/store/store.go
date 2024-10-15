package store

import (
	"fmt"
	"redis-clone/app/internal/store/rdb"
	"time"
)

type data struct {
	value string;
	createdAt uint;
	expireAt uint64;
}
var store = make(map[string]data)

func SetStore(key, val string, args ...uint) {
	timestamp := uint(time.Now().UnixMilli())

	if len(args) > 0 {
		store[key] = data {
			value: val,
			createdAt: timestamp,
			expireAt: uint64(timestamp + args[0]),
		}
	} else {
		store[key] = data {
			value: val,
			createdAt: timestamp,
		}
	}
}

func GetStore(key string) (string, bool) {
    // Check if the key is present in the store
    val, isPresent := store[key]
    fmt.Printf("Val: %s\n", val.value)
    if !isPresent {
        return "", false // Key doesn't exist
    }
    
    // If there is an expiration set and it's expired, remove the key
    if val.expireAt != 0 && uint64(time.Now().UnixMilli()) >= val.expireAt {
        delete(store, key)
        return "", false
    }

    // Otherwise, return the value and true (indicating the key exists and hasn't expired)
    return val.value, true
}

func GetKeys() []string {
	keys := []string{}
	for key, val := range store {
		if val.expireAt > 0 && uint64(time.Now().UnixMilli()) >= val.expireAt {
			delete(store, key)
		} else {
			keys = append(keys, key)
		}
	}
	return keys
}

func LoadRdbInStore(dirPath, fileName string) (bool, error) {
	parsedRdb, err := rdb.ParseRdbFile(dirPath, fileName)

	if err != nil {
		return false, err
	}

	for _, database := range parsedRdb.Databases {
		for key, val := range database.KVMap {
			store[key] = data {
				value: val.Value,
				expireAt: val.ExpireAt,
			}
		}
	}

	return true, nil
}