package scheduler

import (
	"fmt"
	"bytes"
)

type Scheduler interface {
	Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision
}

func ToString(decisions []Decision) string {
	var buffer bytes.Buffer
	for i, d := range decisions {
		buffer.WriteString(fmt.Sprintf("p#%d JobIdx=%d WorkerIdx=%d\n", i, d.JobIdx, d.WorkerIdx))
	}
	return buffer.String()
}

func (d Decision) Equal(x Decision) bool {
	return (d.WorkerIdx == x.WorkerIdx) && (d.JobIdx == x.JobIdx)
}

type MainScheduler struct {
	Scheduler
}

func (m *MainScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs
	d := []Decision{}
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if (j.CPU <= w.CPU) && (j.RAMmb <= w.RAMmb) {
				//add allocation decision to result slice
				d = append(d, Decision{JobIdx: j.Id, WorkerIdx: w.Id})
				//kick allocated worker
				workers = append(workers[:i], workers[i+1:]...)
				break
			}
		}
	}
	return d
}

type GreatScheduler struct {
	Scheduler
	available map[string]Organization
}

func (g *GreatScheduler) init(jobs []ResourceVolume) bool {
	if g.available == nil {
		g.available = make(map[string]Organization)
	}
	for _, j := range jobs {
		if _, ok := g.available[j.Owner.Name]; !ok {
			g.available[j.Owner.Name] = *j.Owner
		}
	}
	return true
}

func (g *GreatScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs adapted
	d := []Decision{}
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if (g.checkAvailability(j)) && (j.CPU <= w.CPU) && (j.RAMmb <= w.RAMmb) {
				if (g.checkQuota(j, w)) {
					//add allocation decision to result slice
					d = append(d, Decision{JobIdx: j.Id, WorkerIdx: w.Id})
					//kick allocated worker
					workers = append(workers[:i], workers[i+1:]...)
					break
				}
			}
		}
	}
	return d
}

func (g *GreatScheduler) checkAvailability(job ResourceVolume) bool {
	return 1.0/float32(len(g.available)) <= job.Owner.Quota.ProjectRatio
}

func (g *GreatScheduler) checkQuota(job ResourceVolume, worker ResourceVolume) bool {
	switch f := job.Owner.Quota.Q.(type) {
	case *Quotum_CpuTimeAbs:
		x := g.available[job.Owner.Name].Quota.GetCpuTimeAbs()

		if !(x >= job.TimePeriod) {
			//re-assign
			g.available[job.Owner.Name] = Organization{job.Owner.Name, &Quotum{
				g.available[job.Owner.Name].Quota.GetProjectRatio(), &Quotum_CpuTimeAbs{x - job.TimePeriod}}}
			return false
		}
	case *Quotum_CpuTimeRatio:
		return true
	case *Quotum_GbAbs:
		return true
	case *Quotum_GbRatio:
		return true
	case nil:
		fmt.Errorf("owner.Quota The field is not set. %T", f)
		return false
	default:
		fmt.Errorf("owner.Quota has unexpected type %T", f)
		return false
	}
	//re-assign due to 
	x := 1.0 / float32(len(g.available))
	g.available[job.Owner.Name] = Organization{job.Owner.Name, &Quotum{
		g.available[job.Owner.Name].Quota.GetProjectRatio() - x, job.Owner.Quota.Q}}
	return true
}
