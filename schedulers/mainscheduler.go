package scheduler

type Scheduler interface {
	Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision
}
type MainScheduler struct {
	Scheduler
}

func (m MainScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	return []Decision{}
}
