package scheduler

import (
	"bytes"
	"fmt"
)

type Scheduler interface {
	Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision
}

func ToString(decisions []Decision) string {
	var buffer bytes.Buffer
	for i, d := range decisions {
		buffer.WriteString(fmt.Sprintf("pair #%s JobIdx %s WorkerIdx %s\n", i, d.JobIdx, d.WorkerIdx))
	}
	return buffer.String()
}

func (d Decision) Equal(x Decision) bool {
	return (d.WorkerIdx == x.WorkerIdx) && (d.JobIdx == x.JobIdx)
}

type MainScheduler struct {
	Scheduler
}

func (m MainScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs
	d := []Decision{}
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if (j.TimePeriod <= w.TimePeriod) && (j.CPU <= w.CPU) && (j.GPU <= w.GPU) && (j.RAMmb <= w.RAMmb) {
				//add allocation decision to result slice
				d = append(d, Decision{JobIdx: j.Id, WorkerIdx: w.Id})
				//kick allocated worker
				workers = append(workers[:i], workers[i+1:]...)
				break
			}
		}
	}
	return d
}
