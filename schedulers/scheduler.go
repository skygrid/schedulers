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
	Available       map[string]Organization
	Counter         map[string]int64
	CpuHoursCounter map[string]float32
	GbCounter       map[string]float32
}

func (g *QuotaScheduler) init() bool {
	if g.Available == nil {
		g.Available = make(map[string]Organization)
	}
	if g.Counter == nil {
		g.Counter = make(map[string]int64)
	}
	if g.CpuHoursCounter == nil {
		g.CpuHoursCounter = make(map[string]float32)
	}
	if g.GbCounter == nil {
		g.GbCounter = make(map[string]float32)
	}
	return true
}

func (g *QuotaScheduler) update(jobs []ResourceVolume) bool {
	for _, j := range jobs {
		g.Available[j.Owner.Name] = *j.Owner
	}
	return true
}

func (g *QuotaScheduler) incrementCounters(job ResourceVolume) {
	g.Counter[job.Owner.GetName()] = g.Counter[job.Owner.GetName()] + 1
	g.CpuHoursCounter[job.Owner.GetName()] = g.CpuHoursCounter[job.Owner.GetName()] + float32(job.TimePeriod)/3600.0
	g.GbCounter[job.Owner.GetName()] = g.GbCounter[job.Owner.GetName()] + job.GetTemporaryStorageNeededGb()
}

func (g *QuotaScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs adapted
	d := []Decision{}
	for _, j := range jobs {
		//first fit
		for i, w := range workers {
			//check availability
			if (g.checkProjectRatio(j)) && (j.CPU <= w.CPU) && (j.RAMmb <= w.RAMmb) && (j.GPU <= w.GPU) {
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

func (g *QuotaScheduler) checkProjectRatio(job ResourceVolume) bool {
	sum := int64(0)
	mul := int64(0)
	for _, v := range g.Counter {
		sum = sum + v
		mul = mul * v
	}
	if (mul == 0) || (float32(g.Counter[job.Owner.GetName()])/float32(sum) <= job.Owner.Quota.GetProjectRatio()) {
		return true
	}
	return false
}

func (g *QuotaScheduler) checkCpuTimeRatio(job ResourceVolume) bool {
	sum := float32(0)
	mul := float32(0)
	for _, v := range g.CpuHoursCounter {
		sum = sum + v
		mul = mul * v
	}
	if (mul == 0) || (float32(g.CpuHoursCounter[job.Owner.GetName()])/float32(sum) <= job.Owner.Quota.GetCpuTimeRatio()) {
		return true
	}
	return false
}

func (g *QuotaScheduler) checkGbRatio(job ResourceVolume) bool {
	sum := float32(0)
	mul := float32(0)
	for _, v := range g.GbCounter {
		sum = sum + v
		mul = mul * v
	}
	if (mul == 0) || (float32(g.GbCounter[job.Owner.GetName()])/float32(sum) <= job.Owner.Quota.GetGbRatio()) {
		return true
	}
	return false
}

func (g *QuotaScheduler) checkQuota(job ResourceVolume, worker ResourceVolume) bool {
	switch f := job.Owner.Quota.Q.(type) {
	case *Quotum_CpuTimeAbs:
		x := g.Available[job.Owner.GetName()].Quota.GetCpuTimeAbs()
		//seconds to hours
		timeSecs := uint64(x * 60 * 60)
		// if quota allows - decrease Available
		if timeSecs >= job.TimePeriod {
			//re-assign
			g.Available[job.Owner.GetName()] = Organization{job.Owner.GetName(), &Quotum{
				g.Available[job.Owner.GetName()].Quota.GetProjectRatio(), &Quotum_CpuTimeAbs{float32(timeSecs-job.GetTimePeriod()) / 3600.0}}}
			g.incrementCounters(job)
		} else {
			return false
		}
	case *Quotum_CpuTimeRatio:
		if g.checkCpuTimeRatio(job) {
			return true
		}
		return false
	case *Quotum_GbAbs:
		x := g.Available[job.Owner.GetName()].Quota.GetGbAbs()
		// if quota allows - decrease Available
		if x >= job.GetTemporaryStorageNeededGb() {
			//re-assign
			g.Available[job.Owner.GetName()] = Organization{job.Owner.GetName(), &Quotum{
				g.Available[job.Owner.GetName()].Quota.GetProjectRatio(), &Quotum_GbAbs{x - job.GetTemporaryStorageNeededGb()}}}
			g.incrementCounters(job)
		} else {
			return false
		}
	case *Quotum_GbRatio:
		if g.checkGbRatio(job) {
			return true
		}
		return false
	case nil:
		fmt.Errorf("owner.Quota The field is not set. %T", f)
		return false
	default:
		fmt.Errorf("owner.Quota has unexpected type %T", f)
		return false
	}
	return true
}
