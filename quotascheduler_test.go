package schedulers

import (
	"testing"
)

func decisionsEqual(a Decision, b Decision) bool {
	if a.JobIdx == b.JobIdx && a.WorkerIdx == b.WorkerIdx {
		return true
	}
	return false
}

func Test_1(t *testing.T) {
	sqs := NewQuotaScheduler()

	// both 50% project weight , 100 CPUhours
	quota1 := Quotum{ProjectRatio: 0.5, Q: &Quotum_CpuHoursAbs{100}}
	quota2 := Quotum{ProjectRatio: 0.5, Q: &Quotum_CpuHoursAbs{100}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 2, TimePeriod: 1000, Owner: &o1, Id: 1}
	job2 := ResourceVolume{CPU: 3, TimePeriod: 900, Owner: &o2, Id: 2}

	//collecting
	jobs1 := []ResourceVolume{job1, job2}

	worker := ResourceVolume{CPU: 3, Id: 10}

	d1 := sqs.scheduleOne(jobs1, worker)
	decision1 := Decision{JobIdx: 1, WorkerIdx: 10}

	if !decisionsEqual(d1, decision1) {
		t.Fail()
	}
	t.Log(sqs.Counter)

	jobs2 := []ResourceVolume{job1, job2}
	d2 := sqs.scheduleOne(jobs2, worker)
	decision2 := Decision{JobIdx: 2, WorkerIdx: 10}

	if !decisionsEqual(d2, decision2) {
		t.Fail()
	}
	t.Log(sqs.Counter)
}

func Test_2(t *testing.T) {
	qs := NewQuotaScheduler()

	// both 50% project weight , 100 CPUhours
	quota1 := Quotum{ProjectRatio: 0.5, Q: &Quotum_CpuHoursAbs{100}}
	quota2 := Quotum{ProjectRatio: 0.5, Q: &Quotum_CpuHoursAbs{100}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 2, TimePeriod: 1000, Owner: &o1, Id: 1}
	job2 := ResourceVolume{CPU: 3, TimePeriod: 900, Owner: &o2, Id: 2}

	//collecting
	jobs1 := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{CPU: 2, Id: 10}
	worker2 := ResourceVolume{CPU: 3, Id: 11}

	workers := []ResourceVolume{worker1, worker2}

	d1 := qs.Schedule(jobs1, workers)
	decisions1 := []Decision{{JobIdx: 1, WorkerIdx: 10}, {JobIdx: 2, WorkerIdx: 11}}

	if !checkDecisionsEqual(d1, decisions1) {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs2 := []ResourceVolume{job1, job2}
	workers2 := []ResourceVolume{worker2}
	d2 := qs.Schedule(jobs2, workers2)
	decision2 := Decision{JobIdx: 1, WorkerIdx: 11}

	if !decisionsEqual(d2[0], decision2) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs3 := []ResourceVolume{job2}
	workers3 := []ResourceVolume{worker1}
	d3 := qs.Schedule(jobs3, workers3)
	decision3 := Decision{}

	if !decisionsEqual(d3[0], decision3) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)
}

func Test_3(t *testing.T) {
	qs := NewQuotaScheduler()

	// both 50% project weight , 100 RAMhours
	quota1 := Quotum{ProjectRatio: 0.5, Q: &Quotum_RamHoursAbs{1000}}
	quota2 := Quotum{ProjectRatio: 0.5, Q: &Quotum_RamHoursAbs{1000}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{RAMmb: 2048, TimePeriod: 1000, Owner: &o1, Id: 1}
	job2 := ResourceVolume{RAMmb: 3096, TimePeriod: 900, Owner: &o2, Id: 2}

	//collecting
	jobs1 := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{RAMmb: 3096, Id: 10}

	d1 := qs.Schedule(jobs1, []ResourceVolume{worker1})
	decision1 := Decision{JobIdx: 1, WorkerIdx: 10}

	if !decisionsEqual(d1[0], decision1) || len(d1) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs2 := []ResourceVolume{job1, job2}
	d2 := qs.Schedule(jobs2, []ResourceVolume{worker1})
	decision2 := Decision{JobIdx: 2, WorkerIdx: 10}

	if !decisionsEqual(d2[0], decision2) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)
}

func Test_4(t *testing.T) {
	qs := NewQuotaScheduler()

	// both 50% project weight , 1500 RAMhours
	quota1 := Quotum{ProjectRatio: 0.5, Q: &Quotum_RamHoursAbs{1500}}
	quota2 := Quotum{ProjectRatio: 0.5, Q: &Quotum_RamHoursAbs{1500}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{RAMmb: 2048, TimePeriod: 1000, Owner: &o1, Id: 1}
	job2 := ResourceVolume{RAMmb: 3096, TimePeriod: 900, Owner: &o2, Id: 2}

	//collecting
	jobs1 := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{RAMmb: 2048, Id: 10}
	worker2 := ResourceVolume{RAMmb: 3096, Id: 11}

	workers := []ResourceVolume{worker1, worker2}

	d1 := qs.Schedule(jobs1, workers)
	decisions1 := []Decision{{JobIdx: 1, WorkerIdx: 10}, {JobIdx: 2, WorkerIdx: 11}}

	if !checkDecisionsEqual(d1, decisions1) {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs2 := []ResourceVolume{job1, job2}
	workers2 := []ResourceVolume{worker2}
	d2 := qs.Schedule(jobs2, workers2)
	decision2 := Decision{JobIdx: 1, WorkerIdx: 11}

	if !decisionsEqual(d2[0], decision2) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs3 := []ResourceVolume{job2}
	workers3 := []ResourceVolume{worker2}
	d3 := qs.Schedule(jobs3, workers3)
	decision3 := Decision{}

	if !decisionsEqual(d3[0], decision3) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)
}

func Test_5(t *testing.T) {
	qs := NewQuotaScheduler()

	// both 50% project weight , 20% and 80% CPUhours
	quota1 := Quotum{ProjectRatio: 0.5, Q: &Quotum_CpuHoursRatio{0.2}}
	quota2 := Quotum{ProjectRatio: 0.5, Q: &Quotum_CpuHoursRatio{CpuHoursRatio: 0.8}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 2, TimePeriod: 1000, Owner: &o1, Id: 1}
	job2 := ResourceVolume{CPU: 3, TimePeriod: 900, Owner: &o2, Id: 2}

	//collecting
	jobs1 := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{CPU: 2, Id: 10}
	worker2 := ResourceVolume{CPU: 3, Id: 11}

	workers := []ResourceVolume{worker1, worker2}

	d1 := qs.Schedule(jobs1, workers)
	decisions1 := []Decision{{JobIdx: 1, WorkerIdx: 10}, {JobIdx: 2, WorkerIdx: 11}}

	if !checkDecisionsEqual(d1, decisions1) {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs2 := []ResourceVolume{job1, job2}
	workers2 := []ResourceVolume{worker2}
	d2 := qs.Schedule(jobs2, workers2)
	decision2 := Decision{JobIdx: 2, WorkerIdx: 11}

	if !decisionsEqual(d2[0], decision2) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs3 := []ResourceVolume{job1, job2}
	workers3 := []ResourceVolume{worker2}
	d3 := qs.Schedule(jobs3, workers3)
	decision3 := Decision{}

	if !decisionsEqual(d3[0], decision3) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)
}

func Test_6(t *testing.T) {
	qs := NewQuotaScheduler()

	//  100% project weight , 2 CPUhours
	quota1 := Quotum{ProjectRatio: 1.0, Q: &Quotum_CpuHoursAbs{CpuHoursAbs: 2}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}

	job1 := ResourceVolume{CPU: 4, TimePeriod: 10000, Owner: &o1, Id: 1}

	//collecting
	jobs1 := []ResourceVolume{job1}

	worker1 := ResourceVolume{CPU: 4, Id: 10}

	d1 := qs.scheduleOne(jobs1, worker1)

	if !decisionsEqual(d1, Decision{}) {
		t.Fail()
	}
	t.Log(qs.Counter)
}

func Test_7(t *testing.T) {
	qs := NewQuotaScheduler()

	// both 50% project weight , 20% and 80% RAMhours
	quota1 := Quotum{ProjectRatio: 0.5, Q: &Quotum_RamHoursRatio{0.2}}
	quota2 := Quotum{ProjectRatio: 0.5, Q: &Quotum_RamHoursRatio{0.8}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{RAMmb: 2048, TimePeriod: 10, Owner: &o1, Id: 1}
	job2 := ResourceVolume{RAMmb: 3096, TimePeriod: 9, Owner: &o2, Id: 2}

	//collecting
	jobs1 := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{RAMmb: 2048, Id: 10}
	worker2 := ResourceVolume{RAMmb: 3096, Id: 11}

	workers := []ResourceVolume{worker1, worker2}

	d1 := qs.Schedule(jobs1, workers)
	decisions1 := []Decision{{JobIdx: 1, WorkerIdx: 10}, {JobIdx: 2, WorkerIdx: 11}}

	if !checkDecisionsEqual(d1, decisions1) {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs2 := []ResourceVolume{job1, job2}
	workers2 := []ResourceVolume{worker2}
	d2 := qs.Schedule(jobs2, workers2)
	decision2 := Decision{JobIdx: 2, WorkerIdx: 11}

	if !decisionsEqual(d2[0], decision2) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)

	jobs3 := []ResourceVolume{job1, job2}
	workers3 := []ResourceVolume{worker2}
	d3 := qs.Schedule(jobs3, workers3)
	decision3 := Decision{}

	if !decisionsEqual(d3[0], decision3) || len(d2) != 1 {
		t.Fail()
	}
	t.Log(qs.Counter)
}

func Test_0(t *testing.T) {
	qs := NewQuotaScheduler()

	//  100% project weight , not initiazed quota
	quota1 := Quotum{ProjectRatio: 1.0}
	o1 := Organization{Name: "SHiP", Quota: &quota1}
	job1 := ResourceVolume{RAMmb: 2048, TimePeriod: 10, Owner: &o1, Id: 1}
	worker1 := ResourceVolume{RAMmb: 2048, Id: 10}

	d1 := qs.scheduleOne([]ResourceVolume{job1}, worker1)
	if !decisionsEqual(d1, Decision{}) {
		t.Fail()
	}
}
