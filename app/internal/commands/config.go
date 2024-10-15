package commands

import (
	"redis-clone/app/internal/resp"
	"redis-clone/app/internal/store"
)

func ConfigGet(key string) (string, error) {
	val, isPresent := store.GetStore("/config/" + key)

	if isPresent {
		response, err := resp.SerializeResp(resp.RespType{
			DataType: resp.Array,
			Array: []*resp.RespType{
				&resp.RespType{
					DataType: resp.BulkString,
					String: key,
				},
				&resp.RespType{
					DataType: resp.BulkString,
					String: val,
				},
			},
		})
		if err != nil {
			return "", err
		}
		return response, nil
	} else {
		return "$-1\r\n", nil
	}
}

func ConfigSet(key, val string) {
	store.SetStore("/config/" + key, val)
}