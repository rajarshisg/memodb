package commands

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"memodb/internal/resp"
)

// HandleCommand function handles the different Redis commands sent by the clients.
func HandleCommand(clientConn net.Conn, buffer []byte) (bool, bool, error) {
	respMsg, err := resp.DeserializeResp(buffer)
	
	// error in deserializing the message
	if err != nil {
		return false, false, err
	}

	var response string
	propagateCommand := false
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
						return false, propagateCommand, err
					}
				}
				case "ECHO": {
					response, err = Echo(arguments)
					if err != nil {
						return false, propagateCommand, err
					}
				}
				case "SET": {
					response, err = Set(arguments)
					propagateCommand = true
					if err != nil {
						return false, propagateCommand, err
					}
				}
				case "GET": {
					response, err = Get(arguments)
					if err != nil {
						return false, propagateCommand, err
					}
				}
				case "KEYS": {
					switch arguments[0] {
						case "*": {
							response, err =	Keys("*")
							if err != nil {
								return false, propagateCommand, err
							}
						}
						default: {
							return false, propagateCommand, fmt.Errorf("unknown command")
						}
					}
				}
				case "CONFIG": {
					switch arguments[0] {
						case "GET": {
							response, err = ConfigGet(arguments[1])
							if err != nil {
								return false, false, err
							}
						}
						default: {
							return false, propagateCommand, fmt.Errorf("unknown command")
						}
					}
				}
				case "INFO": {
					switch arguments[0] {
						case "replication": {
							response, err = InfoReplication()
							if err != nil {
								return false, false, err
							}
						}
						default: {
							return false, propagateCommand, fmt.Errorf("unknown command")
						}
					}
				}
				case "REPLCONF": {
					switch arguments[0] {
						case "listening-port": {
							response, err = ReplConf(clientConn, arguments[1])
							if err != nil {
								return false, propagateCommand, err
							}
						}
						case "capa": {
							response, err = ReplConf(clientConn, "")
							if err != nil {
								return false, propagateCommand, err
							}
						}
						default: {
							return false, propagateCommand, fmt.Errorf("unknown command")
						}
					}
				}
				case "PSYNC": {
					response, err = Psync()
					if err != nil {
						return false, propagateCommand, err
					}

					_, err = clientConn.Write([]byte(response))
					if err != nil {
						fmt.Println("Error while responding to client: ", err.Error())
						return false, propagateCommand, err
					}

					// TODO: We are sending an empty RDB to replica for now, in future this will be replaced by an RDB of the current state
					emptyRdbBytes, _ := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
					clientConn.Write(append([]byte(fmt.Sprintf("$%d\r\n", len(emptyRdbBytes))), emptyRdbBytes...))
					if err != nil {
						fmt.Println("Error while sendign RDB to client: ", err.Error())
						return false, propagateCommand, err
					}
					return true, propagateCommand, nil 
				}
			}
		}
		default: {
			return false,propagateCommand, fmt.Errorf("unknown command")
		}
	}

	_, err = clientConn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error while responding to client: ", err.Error())
		return false, propagateCommand, nil
	}

	return true, propagateCommand, nil
}