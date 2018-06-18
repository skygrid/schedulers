package schedulers

// FCFS Scheduler
type GeneralScheduler struct {
	Scheduler
}

// FCFS scheduling method
func (m *GeneralScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs
	var d []Decision
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if j.CPU <= w.CPU && j.RAMmb <= w.RAMmb && j.GPU <= w.GPU {
				//add allocation decision to result slice
				d = append(d, Decision{JobIdx: j.Id, WorkerIdx: w.Id, CoresNum: w.CPU})
				//kick allocated worker
				workers = append(workers[:i], workers[i+1:]...)
				break
			}
		}
	}
	return d
}
