package scheduler

import (
	"testing"
)

func TestMainScheduler(t *testing.T) {
	m := MainScheduler{}

	// 100% project weight , 100 CPUhours
	quota := Quotum{1.0, &Quotum_CpuTimeAbs{100}}

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
	d_check := []Decision{{JobIdx: 21, WorkerIdx: 21}, {JobIdx: 12, WorkerIdx: 12}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}
