package resp

type DataType string
const (
	String DataType = "+"
	Error DataType = "-"
	Integer DataType = ":"
	BulkString DataType = "$"
	Array DataType = "*"
)

type RespType struct {
	DataType DataType;
	String string;
	Number int;
	Boolean bool;
	Array []*RespType;
}

func IsValidRespDataType(dataType DataType) bool {
	switch dataType {
		case String, Error, Integer, BulkString, Array: {
			return true;
		}
		default: {
			return false;
		}
	}
}