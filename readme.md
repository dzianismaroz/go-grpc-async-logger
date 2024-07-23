In this assignment you will learn how to build a microservice based on the grpc framework

*You cannot use global variables in this task; store what you need in the fields of the structure that lives in the closure*

You will need to implement:

* Generate the necessary code from the proto file
* Microservice base with the ability to stop the server
* ACL - access control from different clients
* Logging system for called methods
* A system for collecting statistics (just counters) for called methods

The microservice will consist of 2 parts:
* Some kind of business logic. In our example, it doesn’t do anything, you just need to call it
* Administration module, where logging and statistics are located

With the first one everything is simple, there is no logic there.

The second one is more interesting. As a rule, in real microservices both logging and statistics work in a single copy, but in our case they will be available via a streaming interface to anyone who connects to the service. This means that 2 logging clients can connect to the service and both will receive a stream of logs. Also, 2 (or more) statistics modules can connect to the service with different intervals for receiving statistics (for example, every 2, 3 and 5 seconds) and it will be sent asynchronously over each interface.

Since asynchron was mentioned, the task will contain goroutines, timers, mutexes, and context with timeouts/completion.

Features of the task:

You need to place the contents of the files service.pb.go and service_grpc.pb.go (which you obtained when generating the proto file) in service.go for uploading as 1 file
You cannot use global variables in this job. Store everything we need in the fields of the structure.

Run tests with go test -v -race

Versions used:
* libprotoc 3.19.3
* protoc-gen-go v1.27.1
* protoc-gen-go-grpc 1.2.0
