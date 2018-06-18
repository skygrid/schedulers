package schedulers

// native interface of Scheduler
type Scheduler interface {
	Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision
}
