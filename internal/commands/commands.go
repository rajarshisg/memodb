package commands

import (
	"fmt"
	"strings"

	"redis-clone/internal/resp"
)

// HandleCommand function handles the different Redis commands sent by the clients.
func HandleCommand(buffer []byte) (string, error) {
	respMsg, err := resp.DeserializeResp(buffer)

	// error in deserializing the message
	if err != nil {
		return "", err
	}

	var response string;
	switch respMsg.DataType {
		case resp.Array: {
			var arrayElems []string
			for i := 0; i < len(respMsg.Array); i++ {
				elemDataType := respMsg.Array[i].DataType

				switch elemDataType {
					case resp.BulkString: {
						arrayElems = append(arrayElems, respMsg.Array[i].String)
					}
				}
			}

			command := strings.ToUpper(arrayElems[0])
			arguments := arrayElems[1:]
			switch command {
				case "PING": {
					response, err = Ping()
					if err != nil {
						return "", err
					}
				}
				case "ECHO": {
					response, err = Echo(arguments)
					if err != nil {
						return "", err
					}
				}
				case "SET": {
					response, err = Set(arguments)
					if err != nil {
						return "", err
					}
				}
				case "GET": {
					response, err = Get(arguments)
					if err != nil {
						return "", err
					}
				}
				case "KEYS": {
					switch arguments[0] {
						case "*": {
							response, err =	Keys("*")
							if err != nil {
								return "", err
							}
						}
						default: {
							return "", fmt.Errorf("unknown command")
						}
					}
				}
				case "CONFIG": {
					switch arguments[0] {
						case "GET": {
							response, err = ConfigGet(arguments[1])
							if err != nil {
								return "", err
							}
						}
						default: {
							return "", fmt.Errorf("unknown command")
						}
					}
				}
				case "INFO": {
					switch arguments[0] {
						case "replication": {
							response, err = InfoReplication()
							if err != nil {
								return "", err
							}
						}
						default: {
							return "", fmt.Errorf("unknown command")
						}
					}
				}
			}
		}
		default: {
			return "", fmt.Errorf("unknown command")
		}
	}

	return response, nil
}