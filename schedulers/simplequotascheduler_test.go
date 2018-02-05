package scheduler

import (
	"testing"
)

func DecisionsEqual(a Decision, b Decision) bool {
	if a.JobIdx == b.JobIdx && a.WorkerIdx == b.WorkerIdx {
		return true
	}
	return false
}

func Test_1(t *testing.T) {
	sqs := SimpleQuotaScheduler{}
	sqs.init()

	// both 50% project weight , 100 CPUhours
	quota1 := Quotum{0.5, &Quotum_CpuHoursAbs{100}}
	quota2 := Quotum{0.5, &Quotum_CpuHoursAbs{100}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 2, TimePeriod: 1000, Owner: &o1, Id: 1}
	job2 := ResourceVolume{CPU: 3, TimePeriod: 900, Owner: &o2, Id: 2}

	//collecting
	jobs := []ResourceVolume{job1, job2}

	worker := ResourceVolume{CPU: 3, Id: 10}

	d1 := sqs.Schedule(jobs, worker)
	decision1 := Decision{JobIdx: 1, WorkerIdx: 10}

	if !DecisionsEqual(d1, decision1) {
		t.Fail()
	}
	t.Log(sqs.Counter)

	d2 := sqs.Schedule(jobs, worker)
	decision2 := Decision{JobIdx: 2, WorkerIdx: 10}

	if !DecisionsEqual(d2, decision2) {
		t.Fail()
	}
	t.Log(sqs.Counter)

}
