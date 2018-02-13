Scheduling library for "[Disneyland](https://github.com/skygrid/disneyland)"
---
[![Go Report Card](https://goreportcard.com/badge/github.com/andreiSaw/schedulers)](https://goreportcard.com/report/github.com/andreiSaw/schedulers)

For contibuting use the [protobuf](https://github.com/golang/protobuf) lib

Project overview
----
Disneyland system comprising of a set of workers(resources) and set of metatasks; 
workers and tasks are characterised by:
* the available/needed number of CPU cores;
* the available/needed RAM memory in megabytes;
* the available/needed GPU memory in megabytes;
* the time, which could be allocated for a task computing / which requested to compute successfully certain task.

All characteristics are assumed to be represented as finite integer numbers. 
For every task planner should match a worker and create a `Decision`. 
And matchmaking result is assumed to be a set of decisions

Every task has its `Owner` and all owners are grouped into project groups, 
each group has their own quota per day, also per resource; 
planner switches on when: 
1) some tasks were just ended; 
2) there is a new task arrived, it’s necessary to compute them and there is free resources available.

The workers capacity almost always correlate with task’s requested characteristics. 
So, any task is assumed to fit any worker.  
It is assumed that workers are coming to scheduler through equal time periods. 
And if a certain task hasn’t matched with arrived worker than the task could fit a another worker, that will come with time

***

Library consists of two general parts: 
the **protobuf library messages description** and
the **implemented Go library**

**Go library** defines two schedulers based on matching (task x resource):

1. `GeneralScheduler` just matching every single task with available resource

2. `QuotaScheduler` matching every single task with available resource due to 
allocated quota restrictions
    + `Quota` is the limit of CPU or memory usage  per single `Organization`, 
    it could be **fixed** and **relative** due to other `Organzation`'s
    Also it indicates weight of `Organzation` tasks 
    (the more weight is, the more `Organzation` tasks will be computed)
    + E.g. to indicate that `Organzation` named 'SHiP' has 100 CPU-hours per week and weight is 50%
      ```
        quota1 := Quotum{0.5, &Quotum_CpuHoursAbs{100}}
        o1 := Organization{Name: "SHiP", Quota: &quota1}
        ```

Generating client-library
----
To use this library you need to install `go get -u github.com/golang/protobuf` 
then add The compiler plugin, protoc-gen-go, will be installed in `$GOBIN`, defaulting to `$GOPATH/bin`. 
It must be in your `$PATH` for the protocol compiler, `protoc`, to find it.

Generate client-library using `./gen.sh` and move the generated `schedmessages.pb.go` to your project folder


Library usage
---
Create two owners (50% weight both and 100 CPU-hours per day) each with one task and only 1 worker for them
```
qs := QuotaScheduler{}
qs.Init()

o1 := Organization{Name: "SHiP", Quota: &Quotum{0.5, &Quotum_CpuHoursAbs{100}}
o2 := Organization{Name: "Monte_carlo", Quota: &Quotum{0.5, &Quotum_CpuHoursAbs{100}}}

job1 := ResourceVolume{CPU: 2, TimePeriod: 1000, Owner: &o1, Id: 1}
job2 := ResourceVolume{CPU: 3, TimePeriod: 900, Owner: &o2, Id: 2}

jobs := []ResourceVolume{job1, job2}

workers := []ResourceVolume{ResourceVolume{CPU: 3, Id: 10}}

d := qs.Schedule(jobs, workers)
```
