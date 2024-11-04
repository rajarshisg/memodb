package commands

import (
	"fmt"

	"memodb/internal/resp"
	"memodb/internal/worker"
)

func InfoReplication() (string, error) {
	return resp.SerializeResp(resp.RespType{
		DataType: resp.BulkString,
		String: fmt.Sprintf("role:%s\nmaster_replid:%s\nmaster_repl_offset:%d", worker.GetWorkerDetails().Role, worker.GetWorkerDetails().Id, 0),
	})
}