package scheduler

import "fmt"

// mini QuotaScheduler
type simpleQuotaScheduler struct {
	Scheduler
	Counter           map[string]int64
	Quotas            map[string]*Quotum
	CpuHoursCounter   map[string]float32
	GbCounter         map[string]float32
	RamMbHoursCounter map[string]float32
}

//update maps
func (sqs *simpleQuotaScheduler) update(job ResourceVolume) {
	name := job.Owner.GetName()
	if _, ok := sqs.Quotas[name]; !ok {
		sqs.Counter[name] = 0
		sqs.Quotas[name] = job.Owner.GetQuota()
		sqs.CpuHoursCounter[name] = 0
		sqs.GbCounter[name] = 0
		sqs.RamMbHoursCounter[name] = 0
	}
}

// check if all owners are in jobs pool
func (sqs *simpleQuotaScheduler) checkProjectsInJobList(jobs []ResourceVolume) bool {
	mul := 1
	flag := false
	for Akey := range sqs.Quotas {
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
		return false
	}
	return true
}

// check whether ProjectQouta exceeded
func (sqs *simpleQuotaScheduler) checkProjectQouta(job ResourceVolume, prFlag bool) bool {
	sum := int64(0)
	mul := int64(1)
	for _, v := range sqs.Counter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(sqs.Counter[job.Owner.GetName()]) / float32(sum)
	if mul == 0 || x <= job.Owner.Quota.GetProjectRatio() || !prFlag {
		return true
	}
	return false
}

// check whether CPU-hours Quota exceeded
func (sqs *simpleQuotaScheduler) checkCpuHoursQouta(job ResourceVolume, prFlag bool) bool {
	sum := float32(0)
	mul := float32(1)
	for _, v := range sqs.CpuHoursCounter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(sqs.CpuHoursCounter[job.Owner.GetName()]) / float32(sum)
	if mul == 0 || x <= job.Owner.Quota.GetCpuHoursRatio() || !prFlag {
		return true
	}
	return false
}

// check whether RAM-hours Quota exceeded
func (sqs *simpleQuotaScheduler) checkRamHoursQouta(job ResourceVolume, prFlag bool) bool {
	sum := float32(0)
	mul := float32(1)
	for _, v := range sqs.RamMbHoursCounter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(sqs.RamMbHoursCounter[job.Owner.GetName()]) / float32(sum)
	if mul == 0 || x <= job.Owner.Quota.GetRamHoursRatio() || !prFlag {
		return true
	}
	return false
}

// detects Quota type and check if it's exceeded on not
func (sqs *simpleQuotaScheduler) checkQuota(job ResourceVolume, prFlag bool) bool {
	name := job.Owner.GetName()
	switch f := job.Owner.Quota.Q.(type) {
	case *Quotum_CpuHoursAbs:
		x := job.GetTimePeriod() * uint64(job.GetCPU()) / 3600.0
		if float32(x) <= sqs.Quotas[name].GetCpuHoursAbs()-sqs.CpuHoursCounter[name] {
			return true
		}
		return false
	case *Quotum_CpuHoursRatio:
		if sqs.checkCpuHoursQouta(job, prFlag) {
			return true
		}
		return false
	case *Quotum_RamHoursAbs:
		x := job.GetTimePeriod() * uint64(job.GetRAMmb()) / 3600.0
		if float32(x) <= sqs.Quotas[name].GetRamHoursAbs()-sqs.RamMbHoursCounter[name] {
			return true
		}
		return false
	case *Quotum_RamHoursRatio:
		if sqs.checkRamHoursQouta(job, prFlag) {
			return true
		}
		return false
	case nil:
		fmt.Printf("owner.Quota The field is not set. %T", f)
		return false
	default:
		fmt.Printf("owner.Quota has unexpected type %T", f)
		return false
	}
}

// increment map counters
func (sqs *simpleQuotaScheduler) incrementCounters(job ResourceVolume) {
	cpuHours := float32(job.GetTimePeriod()*uint64(job.GetCPU())) / 3600.0
	ramMbHours := float32(job.GetTimePeriod()*uint64(job.GetRAMmb())) / 3600.0

	sqs.Counter[job.Owner.GetName()] += 1
	sqs.GbCounter[job.Owner.GetName()] += 1
	sqs.CpuHoursCounter[job.Owner.GetName()] += cpuHours
	sqs.RamMbHoursCounter[job.Owner.GetName()] += ramMbHours
}

// matching one worker with one job from pool
func (sqs *simpleQuotaScheduler) scheduleOne(jobs []ResourceVolume, w ResourceVolume) Decision {
	prFlag := sqs.checkProjectsInJobList(jobs)
	d := Decision{}
	for _, j := range jobs {
		sqs.update(j)
		if j.CPU <= w.CPU && j.RAMmb <= w.RAMmb && j.GPU <= w.GPU {
			if sqs.checkQuota(j, prFlag) && sqs.checkProjectQouta(j, prFlag) {
				sqs.incrementCounters(j)
				d = Decision{JobIdx: j.Id, WorkerIdx: w.Id}
				break
			}
		}
	}
	return d
}

// after scheduling kicks allocated job from pool
func (sqs *simpleQuotaScheduler) kickAllocatedJob(d Decision, jobs []ResourceVolume) []ResourceVolume {
	for jIx, j := range jobs {
		if j.Id == d.JobIdx {
			//kick allocated job
			jobs = append(jobs[:jIx], jobs[jIx+1:]...)
			break
		}
	}
	return jobs
}
