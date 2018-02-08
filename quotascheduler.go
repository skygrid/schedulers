package scheduler

type QuotaScheduler struct {
	SimpleQuotaScheduler
}

func (q *QuotaScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	var d []Decision
	var decision Decision
	for _, worker := range workers {
		decision, jobs = q.ScheduleOne(jobs, worker)
		d = append(d, decision)
	}
	return d
}
