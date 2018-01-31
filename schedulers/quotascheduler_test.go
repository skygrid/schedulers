package scheduler

import (
	"testing"
)

func TestGreatSchedulerCPU_abs(t *testing.T) {
	LOG_SWITCH := false

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// both 50% project weight , 100 CPUhours
	quota1 := Quotum{0.5, &Quotum_CpuTimeAbs{100}}
	quota2 := Quotum{0.5, &Quotum_CpuTimeAbs{100}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 1, RAMmb: 1, TimePeriod: 10, Owner: &o1, Id: 1, TemporaryStorageNeededGb: 0.8}
	job2 := ResourceVolume{CPU: 1, RAMmb: 1, TimePeriod: 90, Owner: &o2, Id: 2, TemporaryStorageNeededGb: 0.2}

	//collecting
	jobs := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 3}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 4}
	//collecting
	workers := []ResourceVolume{worker1, worker2}

	g.update(jobs)

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 3}, {JobIdx: 2, WorkerIdx: 4}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}

func TestGreatSchedulerGB_abs(t *testing.T) {
	LOG_SWITCH := false

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// both 50% project weight , 1 GB
	quota1 := Quotum{0.5, &Quotum_GbAbs{1.0}}
	quota2 := Quotum{0.5, &Quotum_GbAbs{1.0}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}

	job1 := ResourceVolume{CPU: 1, RAMmb: 1, TimePeriod: 10, Owner: &o1, Id: 1, TemporaryStorageNeededGb: 0.8}
	job2 := ResourceVolume{CPU: 1, RAMmb: 1, TimePeriod: 90, Owner: &o2, Id: 2, TemporaryStorageNeededGb: 0.2}

	//collecting
	jobs := []ResourceVolume{job1, job2}

	worker1 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 3}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 4}
	//collecting
	workers := []ResourceVolume{worker1, worker2}

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	//init project overview
	g.update(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 3}, {JobIdx: 2, WorkerIdx: 4}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}

}

func TestGreatSchedulerGB_abs_complex(t *testing.T) {
	LOG_SWITCH := false

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// all 25% project weight
	quota1 := Quotum{0.25, &Quotum_GbAbs{1.0}}
	quota2 := Quotum{0.25, &Quotum_GbAbs{2.0}}
	quota3 := Quotum{0.25, &Quotum_GbAbs{1.0}}
	quota4 := Quotum{0.25, &Quotum_GbAbs{1.0}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}
	o3 := Organization{Name: "Addit1", Quota: &quota3}
	o4 := Organization{Name: "Test2", Quota: &quota4}

	job1 := ResourceVolume{CPU: 1, RAMmb: 1, Owner: &o1, Id: 1, TemporaryStorageNeededGb: 0.2}
	job2 := ResourceVolume{CPU: 1, RAMmb: 1, Owner: &o2, Id: 2, TemporaryStorageNeededGb: 1.8}
	job3 := ResourceVolume{CPU: 2, RAMmb: 1, Owner: &o3, Id: 3, TemporaryStorageNeededGb: 0.8}
	job4 := ResourceVolume{CPU: 1, RAMmb: 1, Owner: &o4, Id: 4, TemporaryStorageNeededGb: 0.2}

	//collecting
	jobs := []ResourceVolume{job1, job2, job3, job4}

	worker1 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 5}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 6}
	worker3 := ResourceVolume{CPU: 1, RAMmb: 2, Id: 7}
	worker4 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 8}
	//collecting
	workers := []ResourceVolume{worker1, worker2, worker3, worker4}

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	g.update(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 5}, {JobIdx: 2, WorkerIdx: 6}, {JobIdx: 3, WorkerIdx: 8}, {JobIdx: 4, WorkerIdx: 7}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}

func TestGreatSchedulerGB_abs_complex2(t *testing.T) {
	LOG_SWITCH := false

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// all 25% project weight
	quota1 := Quotum{0.25, &Quotum_GbAbs{1.0}}
	quota2 := Quotum{0.25, &Quotum_GbAbs{2.0}}
	quota3 := Quotum{0.25, &Quotum_GbAbs{1.0}}
	quota4 := Quotum{0.25, &Quotum_GbAbs{1.0}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}
	o3 := Organization{Name: "Addit1", Quota: &quota3}
	o4 := Organization{Name: "Test2", Quota: &quota4}

	job1 := ResourceVolume{CPU: 1, RAMmb: 1, Owner: &o1, Id: 1, TemporaryStorageNeededGb: 1.2}
	job2 := ResourceVolume{CPU: 1, RAMmb: 1, Owner: &o2, Id: 2, TemporaryStorageNeededGb: 1.8}
	job3 := ResourceVolume{CPU: 2, RAMmb: 1, Owner: &o3, Id: 3, TemporaryStorageNeededGb: 0.8}
	job4 := ResourceVolume{CPU: 1, RAMmb: 1, Owner: &o4, Id: 4, TemporaryStorageNeededGb: 0.2}

	//collecting
	jobs := []ResourceVolume{job1, job2, job3, job4}

	worker1 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 5}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 6}
	worker3 := ResourceVolume{CPU: 1, RAMmb: 2, Id: 7}
	worker4 := ResourceVolume{CPU: 2, RAMmb: 2, Id: 8}
	//collecting
	workers := []ResourceVolume{worker1, worker2, worker3, worker4}

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	g.update(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 2, WorkerIdx: 5}, {JobIdx: 3, WorkerIdx: 6}, {JobIdx: 4, WorkerIdx: 7}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}

func TestQoutaSchedulerGB_ratio_1(t *testing.T) {
	LOG_SWITCH := true

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// all 25% project weight
	quota1 := Quotum{0.5, &Quotum_GbRatio{0.5}}
	quota2 := Quotum{0.5, &Quotum_GbRatio{0.5}}
	//quota3 := Quotum{0.25, &Quotum_GbRatio{0.35}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}
	//o3 := Organization{Name: "Addit1", Quota: &quota3}

	job1 := ResourceVolume{Owner: &o1, Id: 1, TemporaryStorageNeededGb: 1.2}
	job2 := ResourceVolume{Owner: &o2, Id: 2, TemporaryStorageNeededGb: 1.8}
	job3 := ResourceVolume{Owner: &o1, Id: 3, TemporaryStorageNeededGb: 0.8}

	//collecting
	jobs := []ResourceVolume{job1, job2, job3}

	worker1 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 7}
	worker2 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 8}
	worker3 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 9}

	//collecting
	workers := []ResourceVolume{worker1, worker2, worker3}

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	g.update(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 7}, {JobIdx: 2, WorkerIdx: 8}}
	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}

func TestQoutaSchedulerGB_ratio_2(t *testing.T) {
	LOG_SWITCH := true

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// all 25% project weight
	quota1 := Quotum{0.25, &Quotum_GbRatio{0.5}}
	quota2 := Quotum{0.25, &Quotum_GbRatio{0.15}}
	quota3 := Quotum{0.25, &Quotum_GbRatio{0.2}}
	quota4 := Quotum{0.25, &Quotum_GbRatio{0.15}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}
	o3 := Organization{Name: "Addit1", Quota: &quota3}
	o4 := Organization{Name: "Test2", Quota: &quota4}

	job1 := ResourceVolume{Owner: &o1, Id: 1, TemporaryStorageNeededGb: 1.2}
	job2 := ResourceVolume{Owner: &o2, Id: 2, TemporaryStorageNeededGb: 1.8}
	job3 := ResourceVolume{Owner: &o3, Id: 3, TemporaryStorageNeededGb: 0.8}
	job4 := ResourceVolume{Owner: &o4, Id: 4, TemporaryStorageNeededGb: 0.2}
	job5 := ResourceVolume{Owner: &o1, Id: 5, TemporaryStorageNeededGb: 0.2}
	job6 := ResourceVolume{Owner: &o1, Id: 6, TemporaryStorageNeededGb: 0.2}

	//collecting
	jobs := []ResourceVolume{job1, job2, job3, job4, job5, job6}

	worker1 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 7}
	worker2 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 8}
	worker3 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 9}
	worker4 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 10}
	worker5 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 11}
	worker6 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 12}

	//collecting
	workers := []ResourceVolume{worker1, worker2, worker3, worker4, worker5, worker6}

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	g.update(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 7}, {JobIdx: 2, WorkerIdx: 8}, {JobIdx: 3, WorkerIdx: 9},
		{JobIdx: 4, WorkerIdx: 10}, {JobIdx: 5, WorkerIdx: 11}, {JobIdx: 6, WorkerIdx: 12}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}

/*
func TestQoutaSchedulerGB_ratio_complex(t *testing.T) {
	LOG_SWITCH := true

	g := QuotaScheduler{}
	//init project overview
	g.init()

	// all 25% project weight
	quota1 := Quotum{0.25, &Quotum_GbRatio{0.5}}
	quota2 := Quotum{0.25, &Quotum_GbRatio{0.15}}
	quota3 := Quotum{0.25, &Quotum_GbRatio{0.2}}
	quota4 := Quotum{0.25, &Quotum_GbRatio{0.15}}

	o1 := Organization{Name: "SHiP", Quota: &quota1}
	o2 := Organization{Name: "Monte_carlo", Quota: &quota2}
	o3 := Organization{Name: "Addit1", Quota: &quota3}
	o4 := Organization{Name: "Test2", Quota: &quota4}

	job1 := ResourceVolume{Owner: &o1, Id: 1, TemporaryStorageNeededGb: 1.2}
	job2 := ResourceVolume{Owner: &o2, Id: 2, TemporaryStorageNeededGb: 1.8}
	job3 := ResourceVolume{Owner: &o3, Id: 3, TemporaryStorageNeededGb: 0.8}
	job4 := ResourceVolume{Owner: &o4, Id: 4, TemporaryStorageNeededGb: 0.2}
	job5 := ResourceVolume{Owner: &o1, Id: 5, TemporaryStorageNeededGb: 0.2}
	job6 := ResourceVolume{Owner: &o1, Id: 6, TemporaryStorageNeededGb: 0.2}

	//collecting
	jobs := []ResourceVolume{job1, job2, job3, job4, job5, job6}

	worker1 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 7}
	worker2 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 8}
	worker3 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 9}
	worker4 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 10}
	worker5 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 11}
	worker6 := ResourceVolume{TemporaryStorageNeededGb: 10, Id: 12}

	//collecting
	workers := []ResourceVolume{worker1, worker2, worker3, worker4, worker5, worker6}

	g.update(jobs)

	d := g.Schedule(jobs, workers)

	d_check := []Decision{{JobIdx: 1, WorkerIdx: 7}, {JobIdx: 2, WorkerIdx: 8}, {JobIdx: 3, WorkerIdx: 9},
		{JobIdx: 4, WorkerIdx: 10}, {JobIdx: 5, WorkerIdx: 11}, {JobIdx: 6, WorkerIdx: 12}}

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}

	//no update - next scheduling
	jobs = []ResourceVolume{job2, job3, job4}
	workers = []ResourceVolume{worker2, worker3, worker4}

	if LOG_SWITCH {
		t.Log(Logg(jobs, workers))
	}

	d = g.Schedule(jobs, workers)
	d_check = []Decision{{JobIdx: 3, WorkerIdx: 8}, {JobIdx: 4, WorkerIdx: 9}}

	t.Log(ToString(d_check))
	t.Log(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}
*/
