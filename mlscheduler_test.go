package schedulers

import (
	"os"
	"testing"
)

func TestMlcheduler(t *testing.T) {
	m := MlScheduler{}

	job1 := ResourceVolume{Param: 140, Id: 1}
	job2 := ResourceVolume{Param: 30, Id: 2}
	job3 := ResourceVolume{Param: 5, Id: 3}
	job4 := ResourceVolume{Param: 160, Id: 4}

	//collecting
	jobs := []ResourceVolume{job1, job2, job3, job4}

	worker1 := ResourceVolume{CPU: 4, Id: 10}
	worker2 := ResourceVolume{CPU: 4, Id: 20}
	worker3 := ResourceVolume{CPU: 4, Id: 30}

	//collecting
	workers := []ResourceVolume{worker1, worker2, worker3}
	t.Log("Jobs")
	for _, j := range jobs {
		t.Log(j.toString())
	}
	t.Log("Workers")
	for _, w := range workers {
		t.Log(w.toString())
	}

	d := m.Schedule(jobs, workers)
	dCheck := []Decision{{JobIdx: 4, WorkerIdx: 10, CoresNum: 4}, {JobIdx: 1, WorkerIdx: 20, CoresNum: 4},
		{JobIdx: 2, WorkerIdx: 30, CoresNum: 3}, {JobIdx: 3, WorkerIdx: 30, CoresNum: 1}}
	t.Log("Decision")

	t.Log(toString(dCheck))
	t.Log(toString(d))

	if !checkDecisionsEqual(d, dCheck) {
		t.Fail()
	}
}

func TestMlScheduler_DumpData(t *testing.T) {
	m := MlScheduler{}
	_ = os.Mkdir("temp", os.ModePerm)
	names := []string{"temp/intasks1.json", "temp/machines1.json", "temp/maxcores1.json"}
	m.DumpData([]uint64{140, 30, 5, 160}, 3, 4, names)
}

func TestPacking(t *testing.T) {
	m := MlScheduler{}
	m.Packing()
}

func TestMlScheduler_LoadRes(t *testing.T) {
	m := MlScheduler{}
	name := "temp/out1.json"
	_ = m.LoadRes(name)
}
