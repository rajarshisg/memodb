package commands

import (
	"memodb/internal/resp"
	"memodb/internal/store"
)

func Get(arguments []string) (string, error) {
	key := arguments[0]
	val, isPresent := store.GetStore(key)

	if isPresent {
		response, err := resp.SerializeResp(resp.RespType{
			DataType: resp.BulkString,
			String: val,
		})
		if err != nil {
			return "", err
		}
		return response, nil
	} else {
		return "$-1\r\n", nil
	}

	
}