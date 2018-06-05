package schedulers

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func toString(decisions []Decision) string {
	var buffer bytes.Buffer
	for i, d := range decisions {
		buffer.WriteString(fmt.Sprintf("p#%d J%d W%d C%d \t", i, d.JobIdx, d.WorkerIdx, d.CoresNum))
	}
	return buffer.String()
}

func (volume ResourceVolume) toString() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Id=%d CPU=%d Param=%d ", volume.Id, volume.CPU, volume.Param))
	return buffer.String()
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func checkDecisionsEqual(a []Decision, b []Decision) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		b, _ := in_array(x, b)
		if b {
			return true
		}
	}
	return true
}

func TestMainScheduler(t *testing.T) {
	m := GeneralScheduler{}

	// 100% project weight , 100 CPUhours
	quota := Quotum{ProjectRatio: 1.0}

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
	dCheck := []Decision{{JobIdx: 21, WorkerIdx: 21, CoresNum: 2}, {JobIdx: 12, WorkerIdx: 12, CoresNum: 1}}

	t.Log(toString(dCheck))
	t.Log(toString(d))

	if !checkDecisionsEqual(d, dCheck) {
		t.Fail()
	}
}
