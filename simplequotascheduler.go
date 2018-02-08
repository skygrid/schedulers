package scheduler

import "fmt"

type SimpleQuotaScheduler struct {
	Scheduler
	Counter           map[string]int64
	Quotas            map[string]*Quotum
	CpuHoursCounter   map[string]float32
	GbCounter         map[string]float32
	RamMbHoursCounter map[string]float32
}

func (m *SimpleQuotaScheduler) init() {
	if m.Counter == nil {
		m.Counter = make(map[string]int64)
	}
	if m.Quotas == nil {
		m.Quotas = make(map[string]*Quotum)
	}
	if m.CpuHoursCounter == nil {
		m.CpuHoursCounter = make(map[string]float32)
	}
	if m.GbCounter == nil {
		m.GbCounter = make(map[string]float32)
	}
	if m.RamMbHoursCounter == nil {
		m.RamMbHoursCounter = make(map[string]float32)
	}
}
func (m *SimpleQuotaScheduler) update(job ResourceVolume) {
	name := job.Owner.GetName()
	if _, ok := m.Quotas[name]; !ok {
		m.Counter[name] = 0
		m.Quotas[name] = job.Owner.GetQuota()
		m.CpuHoursCounter[name] = 0
		m.GbCounter[name] = 0
		m.RamMbHoursCounter[name] = 0
	}
}

func (g *SimpleQuotaScheduler) checkProjectsInJobList(jobs []ResourceVolume) bool {
	mul := 1
	flag := false
	for Akey := range g.Quotas {
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

func (m *SimpleQuotaScheduler) checkProjectQouta(job ResourceVolume, prFlag bool) bool {
	sum := int64(0)
	mul := int64(1)
	for _, v := range m.Counter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(m.Counter[job.Owner.GetName()]) / float32(sum)
	if mul == 0 || x <= job.Owner.Quota.GetProjectRatio() || !prFlag {
		return true
	}
	return false
}

func (m *SimpleQuotaScheduler) checkCpuHoursQouta(job ResourceVolume, prFlag bool) bool {
	sum := float32(0)
	mul := float32(1)
	for _, v := range m.CpuHoursCounter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(m.CpuHoursCounter[job.Owner.GetName()]) / float32(sum)
	if mul == 0 || x <= job.Owner.Quota.GetCpuHoursRatio() || !prFlag {
		return true
	}
	return false
}

func (m *SimpleQuotaScheduler) checkRamHoursQouta(job ResourceVolume, prFlag bool) bool {
	sum := float32(0)
	mul := float32(1)
	for _, v := range m.RamMbHoursCounter {
		sum = sum + v
		mul = mul * v
	}
	x := float32(m.RamMbHoursCounter[job.Owner.GetName()]) / float32(sum)
	if mul == 0 || x <= job.Owner.Quota.GetRamHoursRatio() || !prFlag {
		return true
	}
	return false
}

func (g *SimpleQuotaScheduler) checkQuota(job ResourceVolume, prFlag bool) bool {
	name := job.Owner.GetName()
	switch f := job.Owner.Quota.Q.(type) {
	case *Quotum_CpuHoursAbs:
		x := job.GetTimePeriod() * uint64(job.GetCPU()) / 3600.0
		if float32(x) <= g.Quotas[name].GetCpuHoursAbs()-g.CpuHoursCounter[name] {
			return true
		}
		return false
	case *Quotum_CpuHoursRatio:
		if g.checkCpuHoursQouta(job, prFlag) {
			return true
		}
		return false
	case *Quotum_RamHoursAbs:
		x := job.GetTimePeriod() * uint64(job.GetRAMmb()) / 3600.0
		if float32(x) <= g.Quotas[name].GetRamHoursAbs()-g.RamMbHoursCounter[name] {
			return true
		}
		return false
	case *Quotum_RamHoursRatio:
		if g.checkRamHoursQouta(job, prFlag) {
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
}

func (m *SimpleQuotaScheduler) incrementCounters(job ResourceVolume) {
	cpuHours := float32(job.GetTimePeriod()*uint64(job.GetCPU())) / 3600.0
	ramMbHours := float32(job.GetTimePeriod()*uint64(job.GetRAMmb())) / 3600.0

	m.Counter[job.Owner.GetName()] += 1
	m.GbCounter[job.Owner.GetName()] += 1
	m.CpuHoursCounter[job.Owner.GetName()] += cpuHours
	m.RamMbHoursCounter[job.Owner.GetName()] += ramMbHours
}

func (m *SimpleQuotaScheduler) ScheduleOne(jobs []ResourceVolume, w ResourceVolume) (decision Decision, editedJobs []ResourceVolume) {
	prFlag := m.checkProjectsInJobList(jobs)
	for jIx, j := range jobs {
		m.update(j)
		if j.CPU <= w.CPU && j.RAMmb <= w.RAMmb && j.GPU <= w.GPU {
			if m.checkQuota(j, prFlag) && m.checkProjectQouta(j, prFlag) {
				m.incrementCounters(j)
				//kick allocated job
				jobs = append(jobs[:jIx], jobs[jIx+1:]...)
				return Decision{JobIdx: j.Id, WorkerIdx: w.Id}, jobs
			}
		}
	}
	return Decision{}, jobs
}
