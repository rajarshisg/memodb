package commands

import (
	"redis-clone/app/internal/resp"
)

// Echo function handles the ECHO command by returning an str consisting of all the arguments passed
func Echo (arguments []string) (string, error) {
	respMsg := resp.RespType {
		DataType: resp.BulkString,
	}
	strSize := 0
	str := ""
	for idx, argument := range arguments {
		
		if idx != len(arguments) - 1 {
			strSize += len(argument) + 1
			str += argument + " "
		} else {
			strSize += len(argument)
			str += argument
		}
	}
	respMsg.String = str
	
	response, err := resp.SerializeResp(respMsg)
	if err != nil {
		return "", err
	}
	return response, nil
}