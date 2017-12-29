package scheduler

import (
	"testing"
	"fmt"
)

func checkDecisionsEqual(a []Decision, b []Decision) bool {
	for i, x := range a {
		if !x.Equal(b[i]) {
			return false
		}
	}
	return true
}
func TestMainScheduler(t *testing.T) {
	m := MainScheduler{}
	o1 := Organization{Name: "SHiP", Quota: 100}
	rv1 := ResourceVolume{CPU: 2, RAMmb: 1, TimePeriod: 40, Owner: &o1, Id: 21}
	rv2 := ResourceVolume{CPU: 1, RAMmb: 2, TimePeriod: 40, Owner: &o1, Id: 12}
	jobs := []ResourceVolume{rv1, rv2}
	worker1 := ResourceVolume{CPU: 1, RAMmb: 2, Id: 12}
	worker2 := ResourceVolume{CPU: 2, RAMmb: 1, Id: 21}

	workers := []ResourceVolume{worker1, worker2}

	d := m.Schedule(jobs, workers)
	d_check := []Decision{{JobIdx: 21, WorkerIdx: 21}, {JobIdx: 12, WorkerIdx: 12}}

	fmt.Println(ToString(d))

	if !checkDecisionsEqual(d, d_check) {
		t.Fail()
	}
}
