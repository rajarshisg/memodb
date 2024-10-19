package worker

type WorkerType struct {
	Id string
	Role string
	Host string
	Port string
	Master_replid string
	Master_repl_offset int
	Connected_slaves int
	SlavePorts []string
}

var worker = new(WorkerType)

func InitWorker(replica bool, workerHost, workerPort, masterHost, masterPort string) (string, error) {
	if !replica {
		_, err := InitMasterWorker(worker, workerHost, workerPort)

		if err != nil {
			return "", err
		}
		return worker.Id, nil
	} else {
		_, err := InitSlaveWorker(worker, workerHost, workerPort, masterHost, masterPort)

		if err != nil {
			return "", err
		}
		return worker.Id, nil
	}
}

func GetWorkerDetails() *WorkerType {
	return worker
}

func UpdateSlaveDetailsForMaster(slavePort string) bool {
	if worker.Role != "master" {
		return false
	}

	worker.SlavePorts = append(worker.SlavePorts, slavePort)
	return true

}