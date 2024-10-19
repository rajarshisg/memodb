package commands

import (
	"fmt"
	"memodb/internal/resp"
	"memodb/internal/worker"
)


func Psync() (string, error) {
	response, err := resp.SerializeResp(resp.RespType{
		DataType: resp.String,
		String: fmt.Sprintf("FULLRESYNC %s %d", worker.GetWorkerDetails().Master_replid, 0),
	})

	if err != nil {
		return "", err
	}

	return response, nil
}