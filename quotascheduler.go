package scheduler

// QuotaScheduler
type QuotaScheduler struct {
	SimpleQuotaScheduler
}

// scheduling method for QuotaScheduler
// matching workers from pool with tasks from the other pool
func (q *QuotaScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	var d []Decision
	var decision Decision
	for _, worker := range workers {
		decision = q.ScheduleOne(jobs, worker)
		d = append(d, decision)
		jobs = q.kickAllocatedJob(decision, jobs)
	}
	return d
}
