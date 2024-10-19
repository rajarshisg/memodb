package commands

import (
	"memodb/internal/resp"
	"memodb/internal/worker"
)


func ReplConf(slavePort string) (string, error) {
	worker.UpdateSlaveDetailsForMaster(slavePort)
	
	response, err := resp.SerializeResp(resp.RespType{
		DataType: resp.String,
		String: "OK",
	})

	if err != nil {
		return "", err
	}

	return response, nil
}