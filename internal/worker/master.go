package worker

import (
	"fmt"

	"github.com/google/uuid"
)

func InitMasterWorker(worker *WorkerType, workerHost, workerPort string) (bool, error) {
	worker.Id = uuid.NewString()
	worker.Role = "master"
	worker.Host = workerHost
	worker.Port = workerPort
	worker.Master_replid = worker.Id
	worker.Master_repl_offset = 0
	worker.Connected_slaves = 0

	return true, nil
}

func PropagateCommand(buffer []byte) {
	if worker.Role != "master" {
		return
	}

	for _, slave := range  worker.Slaves {
		go func(s Slave) {
            _, err := s.connection.Write(buffer)
            if err != nil {
                fmt.Printf("Error writing to slave %s: %v\n", s.port, err)
            }
        }(slave)
	}
}