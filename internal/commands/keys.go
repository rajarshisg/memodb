package commands

import (
	"strings"

	"memodb/internal/resp"
	"memodb/internal/store"
)

func Keys(pattern string) (string, error){
	keys := store.GetKeys()
	respKeysArr := []*resp.RespType{};
	for _, key := range keys {
		if strings.Contains(key, "/config") {
			continue
		}
		respKeysArr = append(respKeysArr, &resp.RespType{
			DataType: resp.String,
			String: key,
		})
	}
	return resp.SerializeResp(resp.RespType{
		DataType: resp.Array,
		Array: respKeysArr,
	})
}