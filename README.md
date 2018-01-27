Scheduling library for "[Disneyland](https://github.com/skygrid/disneyland)"
---
For contibuting use the [protobuf](https://github.com/golang/protobuf) lib

Project overview
----
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
        quota1 := Quotum{0.5, &Quotum_CpuTimeAbs{100}}
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
