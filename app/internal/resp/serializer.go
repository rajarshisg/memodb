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
		default: {
			return "", fmt.Errorf("error occurred during serialization: data type not found")
		}
	}
}