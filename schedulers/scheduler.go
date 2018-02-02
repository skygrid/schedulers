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
	Available         map[string]Organization
	Counter           map[string]int64
	CpuHoursCounter   map[string]float32
	GbCounter         map[string]float32
	RamMbHoursCounter map[string]float32
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
	if g.RamMbHoursCounter == nil {
		g.RamMbHoursCounter = make(map[string]float32)
	}
	return true
}

func (g *QuotaScheduler) update(jobs []ResourceVolume) bool {
	for _, j := range jobs {
		g.Available[j.Owner.GetName()] = *j.Owner
		g.Counter[j.Owner.GetName()] = 0
		g.GbCounter[j.Owner.GetName()] = 0
		g.CpuHoursCounter[j.Owner.GetName()] = 0
		g.RamMbHoursCounter[j.Owner.GetName()] = 0
	}
	return true
}

func (g *QuotaScheduler) incrementCounters(job ResourceVolume) {
	cpuHours := float32(job.GetTimePeriod()*uint64(job.GetCPU())) / 3600.0
	ramMbHours := float32(job.GetTimePeriod()*uint64(job.GetRAMmb())) / 3600.0

	g.Counter[job.Owner.GetName()] += 1
	g.GbCounter[job.Owner.GetName()] += 1
	g.CpuHoursCounter[job.Owner.GetName()] += cpuHours
	g.RamMbHoursCounter[job.Owner.GetName()] += ramMbHours

	switch job.Owner.Quota.Q.(type) {
	case *Quotum_CpuHoursAbs:
		g.Available[job.Owner.GetName()] = Organization{job.Owner.GetName(), &Quotum{
			g.Available[job.Owner.GetName()].Quota.GetProjectRatio(),
			&Quotum_CpuHoursAbs{g.Available[job.Owner.GetName()].Quota.GetCpuHoursAbs() - cpuHours}}}
		break
	case *Quotum_CpuHoursRatio:
		fmt.Println("not implemented")
		break
	case *Quotum_GbAbs:
		fmt.Println("not implemented")
		break
	case *Quotum_GbRatio:
		fmt.Println("not implemented")
		break
	case *Quotum_RamHoursAbs:
		g.Available[job.Owner.GetName()] = Organization{job.Owner.GetName(), &Quotum{
			g.Available[job.Owner.GetName()].Quota.GetProjectRatio(),
			&Quotum_RamHoursAbs{g.Available[job.Owner.GetName()].Quota.GetRamHoursAbs() - ramMbHours}}}
		break
	case *Quotum_RamHoursRatio:
		fmt.Println("not implemented")
		break
	}
}

func (g *QuotaScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	//fcfs adapted
	d := []Decision{}
	scheduledJobs := []ResourceVolume{}
	prFlag := g.checkProjectsInJobList(jobs)
	n := len(jobs)
	allocatedFlag := false
	exitFlag := false
	jIx := 0
	for (jIx < n) && !exitFlag {
		j := jobs[jIx]
		//first fit
		for wIx, w := range workers {
			fmt.Printf("J%d W%d; ", j.Id, w.Id)
			//check availability
			if (j.CPU <= w.CPU) && (j.RAMmb <= w.RAMmb) && (j.GPU <= w.GPU) {
				if prFlag {
					if g.checkProjectRatio(j) && g.checkQuota(j, w) {
						//add allocation decision to result slice
						d = append(d, Decision{JobIdx: j.Id, WorkerIdx: w.Id})
						scheduledJobs = append(scheduledJobs, j)
						//kick allocated worker and job
						workers = append(workers[:wIx], workers[wIx+1:]...)
						jobs = append(jobs[:jIx], jobs[jIx+1:]...)
						//change variable of length
						n--
						//set flag true
						allocatedFlag = true
						//update counters
						g.incrementCounters(j)
						break
					}
				} else {
					//add allocation decision to result slice
					d = append(d, Decision{JobIdx: j.Id, WorkerIdx: w.Id})
					scheduledJobs = append(scheduledJobs, j)
					//kick allocated worker and job
					workers = append(workers[:wIx], workers[wIx+1:]...)
					jobs = append(jobs[:jIx], jobs[jIx+1:]...)
					//change variable of length
					n--
					//set flag true
					allocatedFlag = true
					//update counters
					g.incrementCounters(j)
					break
				}
			}
			//if it was last worker go to next job
			if wIx == len(workers)-1 {
				jIx++
				fmt.Println(jIx)
			}
		}
		//if it was the last job - we could iterate it over again
		if jIx == n-1 {
			if allocatedFlag {
				exitFlag = false
				allocatedFlag = false
				jIx = 0
			} else {
				exitFlag = true
			}
		}
		//if there is no workers - quit
		if len(workers) == 0 {
			exitFlag = true
		}
	}
	return d
}

func (g *QuotaScheduler) checkProjectRatio(job ResourceVolume) bool {
	sum := int64(0)
	mul := int64(1)
	for _, v := range g.Counter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(g.Counter[job.Owner.GetName()]) / float32(sum)
	if (mul == 0) || (x <= job.Owner.Quota.GetProjectRatio()) {
		return true
	}
	return false
}

//TODO: implement
func (g *QuotaScheduler) checkCpuHoursRatio(job ResourceVolume) bool {
	return false
}

//TODO: implement
func (g *QuotaScheduler) checkRamMbHoursRatio(job ResourceVolume) bool {
	return false
}

//TODO: implement
func (g *QuotaScheduler) checkGbRatio(job ResourceVolume) bool {
	return false
}

func (g *QuotaScheduler) checkQuota(job ResourceVolume, worker ResourceVolume) bool {
	switch f := job.Owner.Quota.Q.(type) {
	case *Quotum_CpuHoursAbs:
		x := job.GetTimePeriod() * uint64(job.GetCPU()) / 3600.0
		if float32(x) < g.Available[job.Owner.GetName()].Quota.GetCpuHoursAbs() {
			return true
		}
		return false
	case *Quotum_CpuHoursRatio:
		fmt.Println("not implemented")
		return false
	case *Quotum_GbAbs:
		fmt.Println("not implemented")
		return false
	case *Quotum_GbRatio:
		fmt.Println("not implemented")
		return false
	case *Quotum_RamHoursAbs:
		x := job.GetTimePeriod() * uint64(job.GetRAMmb()) / 3600.0
		if float32(x) < g.Available[job.Owner.GetName()].Quota.GetRamHoursAbs() {
			return true
		}
		return false
	case *Quotum_RamHoursRatio:
		fmt.Println("not implemented")
		return false
	case nil:
		fmt.Errorf("owner.Quota The field is not set. %T", f)
		return false
	default:
		fmt.Errorf("owner.Quota has unexpected type %T", f)
		return false
	}
	return false
}

func (g *QuotaScheduler) checkProjectsInJobList(jobs []ResourceVolume) bool {
	mul := 1
	retVal := false
	flag := false
	for Akey, _ := range g.Available {
		for _, j := range jobs {
			if j.Owner.GetName() == Akey {
				flag = true
				break
			}
		}
		if !flag {
			mul = 0
		}
		flag = false
	}
	if mul == 0 {
		retVal = false
	} else {
		retVal = true
	}
	return retVal
}
