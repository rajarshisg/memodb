package commands

import (
	"net"

	"memodb/internal/resp"
	"memodb/internal/worker"
)


func ReplConf(clientConn net.Conn, slavePort string) (string, error) {
	if slavePort != "" {
		worker.UpdateSlaveDetailsForMaster(clientConn, slavePort)
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