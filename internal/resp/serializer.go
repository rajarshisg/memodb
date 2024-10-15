package resp

import (
	"fmt"
)

// SerializeResp converts a RespType object into a valid Redis Serialization Protocol (RESP) string.
func SerializeResp(resp RespType) (string, error) {
	dataType := resp.DataType;

	switch dataType {
		case String: {
			return "+" + resp.String + "\r\n", nil
		}
		case BulkString: {
			return "$" + fmt.Sprint(len(resp.String)) + "\r\n" + resp.String + "\r\n", nil
		}
		case Array: {
			size := len(resp.Array)
			response := "*" + fmt.Sprint(size) + "\r\n"

			for _, val := range resp.Array {
				currResponse, err := SerializeResp(*val)
				if err != nil {
					return "", err
				}
				response += currResponse
			}

			return response, nil
		}
		default: {
			return "", fmt.Errorf("error occurred during serialization: data type not found")
		}
	}
}