Scheduling library for "[Disneyland](https://github.com/skygrid/disneyland)"
---
For contibuting use https://github.com/golang/protobuf lib

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

Library installing
----
To use this library you need to install `go get -u github.com/golang/protobuf` then add The compiler plugin, protoc-gen-go, will be installed in `$GOBIN`, defaulting to `$GOPATH/bin`. It must be in your `$PATH` for the protocol compiler, `protoc`, to find it.

Generate client-library using `./gen.sh` and move the generated `schedmessages.pb.go` to your project folder


Library usage
---


Project overview
----
Library defines the protocol buffer description for 