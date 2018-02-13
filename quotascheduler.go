package scheduler

// QuotaScheduler
type QuotaScheduler struct {
	simpleQuotaScheduler
}

// scheduling method for QuotaScheduler
// matching workers from pool with tasks from the other pool
func (q *QuotaScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	var d []Decision
	var decision Decision
	for _, worker := range workers {
		decision = q.scheduleOne(jobs, worker)
		d = append(d, decision)
		jobs = q.kickAllocatedJob(decision, jobs)
	}
	return d
}

// NewQuotaScheduler constructor
func NewQuotaScheduler() *QuotaScheduler {
	sqs := QuotaScheduler{}
	if sqs.Counter == nil {
		sqs.Counter = make(map[string]int64)
	}
	if sqs.Quotas == nil {
		sqs.Quotas = make(map[string]*Quotum)
	}
	if sqs.CpuHoursCounter == nil {
		sqs.CpuHoursCounter = make(map[string]float32)
	}
	if sqs.GbCounter == nil {
		sqs.GbCounter = make(map[string]float32)
	}
	if sqs.RamMbHoursCounter == nil {
		sqs.RamMbHoursCounter = make(map[string]float32)
	}
	return &sqs
}
