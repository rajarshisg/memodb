package resp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DeserializeResp function takes in an input buffer and converts the Redis Serialization Protocol (RESP) message into a standard RespType object.
func DeserializeResp(inputBytes []byte) (*RespType, error) {
	commands := strings.Split(strings.TrimSpace(string(inputBytes)), "\r\n")
	if len(commands) < 1 {
		return nil, errors.New("deserialization error: missing RESP Data Type")
	}

	dataType := DataType(commands[0][:1])
	if !IsValidRespDataType(dataType) {
		return nil, errors.New("deserialization error: malformed RESP data type")
	}

	deserializedMessage := &RespType{DataType: dataType}
	switch dataType {
		case String, Error: {
			if len(commands[0]) < 2 {
				return nil, fmt.Errorf("deserialization error: found RESP data type %s but no data found", dataType)
			}
			deserializedMessage.String = commands[0][1:]
		}
		case BulkString: {
			if len(commands) < 2 {
				return nil, fmt.Errorf("deserialization error: found RESP data type %s but no bulk string found", dataType)
			}
			deserializedMessage.String = commands[1]
		}
		case Integer: {
			if len(commands[0]) < 2 {
				return nil, fmt.Errorf("deserialization error: found RESP data type %s but no integer value found", dataType)
			}
			num, err := strconv.ParseInt(commands[0][1:], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("deserialization error: found RESP data type %s but not an int value", dataType)
			}
			deserializedMessage.Number = int(num)
		}
		case Array: {
			if len(commands) < 2 {
				return nil, fmt.Errorf("deserialization error: found RESP data type %s but no array size found", dataType)
			}
			size, err := strconv.Atoi(commands[0][1:])
			if err != nil {
				return nil, fmt.Errorf("deserialization error: found RESP data type %s but not a valid size", dataType)
			}
			var respArr []*RespType
			count := 0
			idx := 1
			for {
				currElemDataType := DataType(commands[idx][0])

				switch currElemDataType {
					case String, Error, Integer: {
						elem, err := DeserializeResp([]byte(commands[idx]))
						if err != nil {
							return nil, fmt.Errorf("deserialization error: could not process Array element %d", idx)
						}
						respArr = append(respArr, elem)
						idx++
						count++
					}
					case BulkString: {
						elem, err := DeserializeResp([]byte(strings.Join([]string{commands[idx], commands[idx + 1]}, "\r\n")))
						if err != nil {
							return nil, fmt.Errorf("deserialization error: could not process Array element %d", idx)
						}
						respArr = append(respArr, elem)
						idx += 2
						count++
					}
					case Array: {
						// TODO: Need to handle nested arrays
						count++
					}
				}

				if count == size {
					break
				}
			}
			deserializedMessage.Array = respArr
		}
	}

	return deserializedMessage, nil
}
