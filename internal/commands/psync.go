package commands

import (
	"fmt"

	"memodb/internal/resp"
)


func Psync() (string, error) {
	response, err := resp.SerializeResp(resp.RespType{
		DataType: resp.String,
		String: fmt.Sprintf("FULLRESYNC %s %d", "abc", 0),
	})

	if err != nil {
		return "", err
	}

	return response, nil
}