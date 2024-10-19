package worker

import "github.com/google/uuid"

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