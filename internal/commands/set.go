package commands

import (
	"fmt"
	"strconv"

	"redis-clone/internal/resp"
	"redis-clone/internal/store"
)

func Set(arguments []string) (string, error) {
	key := arguments[0]
	val := arguments[1]
	if len(arguments) >= 4 {
		expiryTimeInMilliSeconds, err := strconv.Atoi(arguments[3])
		if err != nil {
			return "", fmt.Errorf("expiry not a valid number in SET command")
		}
		store.SetStore(key, val, uint(expiryTimeInMilliSeconds))
	} else {
		store.SetStore(key, val)
	}

	response, err := resp.SerializeResp(resp.RespType{
		DataType: resp.String,
		String: "OK",
	})
	if err != nil {
		return "", err
	}

	return response, nil
}