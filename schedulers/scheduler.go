package scheduler

import (
	"fmt"
)

type Scheduler interface {
	Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision
}

type GeneralScheduler struct {
	Scheduler
}

func (m *GeneralScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs
	d := []Decision{}
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if (j.CPU <= w.CPU) && (j.RAMmb <= w.RAMmb) && (j.GPU <= w.GPU) {
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

type QuotaScheduler struct {
	Scheduler
	available map[string]Organization
}

func (g *QuotaScheduler) init(jobs []ResourceVolume) bool {
	if g.available == nil {
		g.available = make(map[string]Organization)
	}
	for _, j := range jobs {
		g.available[j.Owner.Name] = *j.Owner
	}
	return true
}

func (g *QuotaScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs adapted
	d := []Decision{}
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if (g.checkAvailability(j)) && (j.CPU <= w.CPU) && (j.RAMmb <= w.RAMmb) && (j.GPU <= w.GPU) {
				if g.checkQuota(j, w) {
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

func (g *QuotaScheduler) checkAvailability(job ResourceVolume) bool {
	return 1.0/float32(len(g.available)) <= job.Owner.Quota.ProjectRatio
}

func (g *QuotaScheduler) checkQuota(job ResourceVolume, worker ResourceVolume) bool {
	switch f := job.Owner.Quota.Q.(type) {
	case *Quotum_CpuTimeAbs:
		x := g.available[job.Owner.Name].Quota.GetCpuTimeAbs()
		//seconds to hours
		timeSecs := uint64(x * 60 * 60)
		// if quota allows - decrease available
		if timeSecs >= job.TimePeriod {
			//re-assign
			g.available[job.Owner.Name] = Organization{job.Owner.Name, &Quotum{
				g.available[job.Owner.Name].Quota.GetProjectRatio(), &Quotum_CpuTimeAbs{float32(timeSecs-job.TimePeriod) / 3600.0}}}
		} else {
			return false
		}
	case *Quotum_CpuTimeRatio:
		return true
	case *Quotum_GbAbs:
		x := g.available[job.Owner.Name].Quota.GetGbAbs()
		// if quota allows - decrease available
		if x >= job.TemporaryStorageNeededGb {
			//re-assign
			g.available[job.Owner.Name] = Organization{job.Owner.Name, &Quotum{
				g.available[job.Owner.Name].Quota.GetProjectRatio(), &Quotum_GbAbs{x - job.TemporaryStorageNeededGb}}}
		} else {
			return false
		}
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
