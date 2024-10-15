package commands

import "redis-clone/internal/resp"

// Ping function handles the PING command by responding with a PONG.
func Ping() (string, error) {
	response, err := resp.SerializeResp(resp.RespType{
		DataType: resp.String,
		String: "PONG",
	})

	if err != nil {
		return "", err
	}

	return response, nil
}