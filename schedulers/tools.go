package scheduler

import (
	"bytes"
	"fmt"
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
	buffer.WriteString(fmt.Sprintf("Id=%d CPU=%d GPU=%d RAM=%d Time=%d Gb=%f Owner %s ", volume.Id, volume.CPU, volume.GPU, volume.RAMmb, volume.TimePeriod, volume.TemporaryStorageNeededGb, volume.Owner))
	return buffer.String()
}

func (d Decision) Equal(x Decision) bool {
	return (d.WorkerIdx == x.WorkerIdx) && (d.JobIdx == x.JobIdx)
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
