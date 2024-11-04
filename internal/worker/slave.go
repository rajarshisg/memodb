package worker

import (
	"fmt"
	"net"
	"strings"

	"memodb/internal/resp"
	"memodb/internal/tcp"

	"github.com/google/uuid"
)

func InitSlaveWorker(worker *WorkerType, workerHost, workerPort, masterHost, masterPort string) (bool, error) {
	worker.Id = uuid.NewString()
	worker.Role = "slave"
	worker.Host = workerHost
	worker.Port = workerPort
	worker.Master_replid = worker.Id
	worker.Master_repl_offset = 0
	worker.Connected_slaves = 0

	isHandshakeSuccess, err := masterHandshake(masterHost, masterPort, workerPort)

	if err != nil {
		return false, err
	}
	if !isHandshakeSuccess {
		return false, fmt.Errorf("could not connect to master")
	}

	return true, nil
}

func masterHandshake(masterHost, masterPort, workerPort string) (bool, error) {
	conn, err := net.Dial("tcp", masterHost + ":" + masterPort)

	if err != nil {
		return false, err
	}

	pingHandshakeSuccess, err := pingHandshake(conn)
	if !pingHandshakeSuccess || err != nil {
		if err == nil {
			return false, fmt.Errorf("error occurred while performing PING handshake with master")
		}
		return false, err
	}

	replConfigHandshakeSuccess, err := replConfigHandshake(conn, fmt.Sprintf("listening-port %s", workerPort))
	if !replConfigHandshakeSuccess || err != nil {
		if err == nil {
			return false, fmt.Errorf("error occurred while performing REPLCONFIG handshake with master")
		}
		return false, err
	}
	replConfigHandshakeSuccess, err = replConfigHandshake(conn, "capa psync2")
	if !replConfigHandshakeSuccess || err != nil {
		if err == nil {
			return false, fmt.Errorf("error occurred while performing REPLCONF handshake with master")
		}
		return false, err
	}
	psyncHandshakeSuccess, err := psyncHandshake(conn)
	if !psyncHandshakeSuccess || err != nil {
		if err == nil {
			return false, fmt.Errorf("error occurred while performing PSYNC handshake with master")
		}
		return false, err
	}

	tcp.Connections = append(tcp.Connections, tcp.TCPConnection{
		Connection: conn,
		Initialized: false,
	})

	return true, nil
}

func pingHandshake(conn net.Conn) (bool, error) {
	pingCommand := resp.RespType{
			DataType: resp.Array,
			Array: []*resp.RespType{{
					DataType: resp.BulkString,
					String: "PING",
				}},
		}
	pingCommandSerialized, _ := resp.SerializeResp(pingCommand)
	
	_, err := conn.Write([]byte(pingCommandSerialized))
	if err != nil {
		return false, err
	}

	responseBytes := make([]byte, 1024)
	_, err = conn.Read(responseBytes)
	if err != nil {
		return false, err
	}

	deserializedPingCommandResp, err := resp.DeserializeResp(responseBytes)
	
	if err != nil {
		return false, err
	}

	if deserializedPingCommandResp.String != "PONG" {
		return false, fmt.Errorf("error in receiving ping response from master")
	}

	return true, nil
}

func replConfigHandshake(conn net.Conn, command string) (bool, error) {
	respArray := []*resp.RespType{{
					DataType: resp.BulkString,
					String: "REPLCONF",
				}}
	for _, argument := range strings.Split(command, " ") {
		respArray = append(respArray, &resp.RespType{
					DataType: resp.BulkString,
					String: argument,
				})
	}
	replConfigCommand := resp.RespType{
			DataType: resp.Array,
			Array: respArray,
		}
	replConfigCommandSerialized, _ := resp.SerializeResp(replConfigCommand)
	_, err := conn.Write([]byte(replConfigCommandSerialized))
	if err != nil {
		return false, err
	}

	responseBytes := make([]byte, 1024)
	_, err = conn.Read(responseBytes)
	if err != nil {
		return false, err
	}

	deserializedReplConfigCommandResp, err := resp.DeserializeResp(responseBytes)
	
	if err != nil {
		return false, err
	}

	if deserializedReplConfigCommandResp.String != "OK" {
		return false, fmt.Errorf("error in receiving ping response from master")
	}

	return true, nil
}

func psyncHandshake(conn net.Conn) (bool, error) {
	psyncCommand := resp.RespType{
			DataType: resp.Array,
			Array: []*resp.RespType{{
					DataType: resp.BulkString,
					String: "PSYNC",
				}, {
					DataType: resp.BulkString,
					String: "?",
				}, {
					DataType: resp.BulkString,
					String: "-1",
				}},
		}
	replConfigCommandSerialized, _ := resp.SerializeResp(psyncCommand)
	_, err := conn.Write([]byte(replConfigCommandSerialized))
	if err != nil {
		return false, err
	}

	responseBytes := make([]byte, 1024)
	_, err = conn.Read(responseBytes)
	if err != nil {
		return false, err
	}

	// deserializedReplConfigCommandResp, err := resp.DeserializeResp(responseBytes)
	
	// if err != nil {
	// 	return false, err
	// }

	// if deserializedReplConfigCommandResp.String != "OK" {
	// 	return false, fmt.Errorf("error in receiving ping response from master")
	// }

	return true, nil
}