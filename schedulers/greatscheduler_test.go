package scheduler

import (
	"testing"
)

func TestGreatScheduler(t *testing.T) {
	g := GreatScheduler{}

	// both 50% project weight , 100 CPUhours
	quota1 := Quotum{0.5, &Quotum_CpuTimeAbs{100}}
	quota2 := Quotum{0.5, &Quotum_CpuTimeAbs{100}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 1, RAMmb: 1, TimePeriod: 10, Owner: &o1, Id: 1, TemporaryStorageNeededGB: 1}
	job2 := ResourceVolume{CPU: 1, RAMmb: 1, TimePeriod: 90, Owner: &o2, Id: 2, TemporaryStorageNeededGB: 1}
	//collecting
	jobs := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 3}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 4}
	//collecting
	workers := []ResourceVolume{worker1, worker2}

	//init project overview
	g.init(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 3}, {JobIdx: 2, WorkerIdx: 4}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}
