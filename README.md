Scheduling library for "[Disneyland](https://github.com/skygrid/disneyland)"
---
For contibuting use the [protobuf](https://github.com/golang/protobuf) lib

Project overview
----
There are users grouped into project groups, 
each group has their own quota per day, also per resource; 
planner switches on when: 
1) some tasks were just ended; 
2) there is a new task arrived, it’s necessary to compute them and there is free resources available.

Library consists of two general parts: 
the **protobuf library messages description** and
the **implemented Go library**

**Go library** defines two schedulers based on matching (job x resource):

1. `GeneralScheduler` just matching every single job with available resource

2. `QuotaScheduler` matching every single job with available resource due to 
allocated quota restrictions
    + `Quota` is the limit of CPU or memory usage  per single `Organization`, 
    it could be **fixed** and **relative** due to other `Organzation`'s
    Also it indicates weight of `Organzation` jobs 
    (the more weight is, the more `Organzation` tasks will be computed)
    + E.g. to indicate that `Organzation` named 'SHiP' has 100 CPU-hours per week and weight is 50%
      ```
        quota1 := Quotum{0.5, &Quotum_CpuHoursAbs{100}}
        o1 := Organization{Name: "SHiP", Quota: &quota1}
        ```

Project structure
---
```$xslt
/schedulers
    generalscheduler_test.go
    quotascheduler_test.go
    scheduler.go
gen.sh
run_tests.sh
schedmessages.proto
```

Generating client-library
----
To use this library you need to install `go get -u github.com/golang/protobuf` 
then add The compiler plugin, protoc-gen-go, will be installed in `$GOBIN`, defaulting to `$GOPATH/bin`. 
It must be in your `$PATH` for the protocol compiler, `protoc`, to find it.

Generate client-library using `./gen.sh` and move the generated `schedmessages.pb.go` to your project folder


Library usage
---
Create two projects (50% weight both) and one job for each and only 1 worker 
```
qs := QuotaScheduler{}
qs.init()

o1 := Organization{Name: "SHiP", Quota: &Quotum{0.5, &Quotum_CpuHoursAbs{100}}
o2 := Organization{Name: "Monte_carlo", Quota: &Quotum{0.5, &Quotum_CpuHoursAbs{100}}}

job1 := ResourceVolume{CPU: 2, TimePeriod: 1000, Owner: &o1, Id: 1}
job2 := ResourceVolume{CPU: 3, TimePeriod: 900, Owner: &o2, Id: 2}

jobs := []ResourceVolume{job1, job2}

worker := []ResourceVolume{ResourceVolume{CPU: 3, Id: 10}}

d := qs.Schedule(jobs, worker)
```
