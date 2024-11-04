package worker

import (
	"memodb/internal/store"
	"net"
)
type Slave struct {
	port string
	connection net.Conn
}
type WorkerType struct {
	Id string
	Role string
	Host string
	Port string
	Master_replid string
	Master_repl_offset int
	Connected_slaves int
	Slaves []Slave
}

var worker = new(WorkerType)

func InitWorker(replica bool, workerHost, workerPort, masterHost, masterPort string) (string, error) {
	if !replica {
		_, err := InitMasterWorker(worker, workerHost, workerPort)

		if err != nil {
			return "", err
		}
		
	} else {
		_, err := InitSlaveWorker(worker, workerHost, workerPort, masterHost, masterPort)

		if err != nil {
			return "", err
		}
	}
	store.SetStore("internal/worker/role", worker.Role)
	store.SetStore("internal/worker/id", worker.Id)
	return worker.Id, nil
}

func GetWorkerDetails() *WorkerType {
	return worker
}

func UpdateSlaveDetailsForMaster(clientCon net.Conn, port string) bool {
	if worker.Role != "master" {
		return false
	}

	worker.Slaves = append(worker.Slaves, Slave{
		port: port,
		connection: clientCon,
	})
	return true
}