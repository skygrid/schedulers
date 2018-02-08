package scheduler

import (
	"bytes"
	"fmt"
	"testing"
)

func ToString(decisions []Decision) string {
	var buffer bytes.Buffer
	for i, d := range decisions {
		buffer.WriteString(fmt.Sprintf("p#%d J%d W%d ", i, d.JobIdx, d.WorkerIdx))
	}
	return buffer.String()
}

func (volume ResourceVolume) ToString() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Id=%d CPU=%d GPU=%d RAM=%d Time=%d Owner %s ", volume.Id, volume.CPU, volume.GPU, volume.RAMmb, volume.TimePeriod, volume.Owner))
	return buffer.String()
}

func (d Decision) Equal(x Decision) bool {
	return d.WorkerIdx == x.WorkerIdx && d.JobIdx == x.JobIdx
}

func checkDecisionsEqual(a []Decision, b []Decision) bool {
	for i, x := range a {
		if !x.Equal(b[i]) {
			return false
		}
	}
	if len(a) != len(b) {
		return false
	}
	return true
}

func Logg(jobs []ResourceVolume, workers []ResourceVolume) string {
	var buffer bytes.Buffer
	buffer.WriteString("\n")

	for _, j := range jobs {
		buffer.WriteString(fmt.Sprintf("%s\n", j.ToString()))
	}
	for _, w := range workers {
		buffer.WriteString(fmt.Sprintf("%s\n", w.ToString()))
	}
	return buffer.String()
}

func TestMainScheduler(t *testing.T) {
	m := GeneralScheduler{}

	// 100% project weight , 100 CPUhours
	quota := Quotum{1.0, &Quotum_CpuHoursAbs{100}}

	o1 := Organization{Name: "SHiP", Quota: &quota}

	job1 := ResourceVolume{CPU: 2, RAMmb: 1, TimePeriod: 40, Owner: &o1, Id: 21}
	job2 := ResourceVolume{CPU: 1, RAMmb: 2, TimePeriod: 40, Owner: &o1, Id: 12}
	//collecting
	jobs := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{CPU: 1, RAMmb: 2, Id: 12}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 1, Id: 21}
	//collecting
	workers := []ResourceVolume{worker1, worker2}

	d := m.Schedule(jobs, workers)
	dCheck := []Decision{{JobIdx: 21, WorkerIdx: 21}, {JobIdx: 12, WorkerIdx: 12}}

	t.Log(ToString(dCheck))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, dCheck) {
		t.Fail()
	}
}
