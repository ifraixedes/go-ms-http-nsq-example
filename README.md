# Go MS HTTP NSQ Example

This implementation was created to fulfill the requirements that you can find in the [requirements.md file](requirements.md), in around 30 hours and with the purpose of containing good practices like testing, decent git commit history and a decent source documentation.

## Setup

At least Go 1.11 is needed for compiling, running and/or building the binaries.

Each _Makefile_ `all` target build the binary of each microservice. Each binary takes some command line flags for parametrizing external systems, server addresses, etc., you can check them using the `-h` command line flag.

As commented in the same general _Makefile_, there is a `e2e` target which do some requests to the gateway to get some responses; this target works if the current docker-compose is running and the 3 binaries are running with the default configuration, considering that the gateway binary requires a configuration file and the one, in the repository, is being used. The order of the execution of these 3 binaries must be:

1. `zombie-driver/http/cmd/server`
2. `driver-location/nsqhttp/cmd/server`
3. `gateway/cmd/server`

The order is in this way to guarantee that the _NSQ topic_ is created before the consumer runs; although NSQ consumer should be autonomous finding the _topic_ through the _NSQ Looupd service_, it isn't a case which has been verified nor dug on it because of the time dedicated to the test and the milestone which was marked.

## Considerations

The commit history is clean and the commits comments contain good explanations on what's done on each commit; I encourage you to read it if you have doubts or you want to have a general idea why and what decisions have been done during the implementation.

The sources are documented using the _go doc_ comment format.

The code is tested, some parts are tested with integration test because some of the internal parts are already tested and because they verify the complete chain of the involved parts. There are some mocks, but only for the service interfaces (see the _implementation details_ sections), there isn't any mock for any external system (e.g. Redis), because it's preferable to test it with the real ones, moreover nowadays, that we can run them in containers.

You can run the tests with the `test` Makefile targets, without changing anything if you run the docker-compose as its, nonetheless if you run the involved external systems in other addresses, you can pass different addresses to the test using Makefile variables.

## Implementation details

The Driver Location and Zombie Driver services are defined with a general _service interface_. This interface define the functionality of the service and each implementation must satisfy it.

An implementation of a service is usually a package with the name of the _main external system_ on which it relies the business logic; for example in the Driver Location Service, the implementation which relies on Redis is inside of subpackage named `redis`.

A part of the implementations, the services require some transport in order of exposing remotely its functionality (you know not being used as static library). The transport is a package whose name is the system or protocol which is used to expose the service. Each transport has a server part, which is a binary which runs the service behinds a server using such system/protocol, but also, it has a client part for easing to use the service from Go sources (the client satisfy the _service interface_, which the service public contract).

It's worth to comment that the Zombie Driver service set of rules is configurable.

The Gateway  is a bit different than the other 2, in the sense that, it doesn't have a interface; this is mostly because of the role of _API gateway_ that it has. The package exports the configuration (which loaded from the YAML) and some functions to configure an HTTP server; a subpackage contains the main which allows to build the binary.

## Improvements

Software isn't never finish, so there are always room to improve, however this is a short list of clear things to be improved:

* If you look for `TODO:` comment in the sources, you'll see place where the code must be improved, besides some test cases has been skipped due the marked milestone to have the implementation. I don't usually use `TODO` comments without being tracked in an issue tracker, but this is not a real project, so they have been left without being tracked.
* The drivers distances are calculated lineally on the globe of the earth, so it has to be changed to be calculated on the actual path that they follow, however that requires to use some map library.
* Metrics and tracing has not been added because I didn't have the time for them;  the concept that I bore in mind was to use [OpenCensus](https://opencensus.io/) [Go package](https://opencensus.io/quickstart/go/) and its [_HTTP_ plugin](https://godoc.org/go.opencensus.io/plugin/ochttp), which is as easy to add them on each service, just doing an implementation which satisfies the _service interface_ and wrap the business logic implementation of the service (i.e modeling the system in layers).
* Obviously circuit breaker should also be added, however I did not have time to actual take a look into it for giving instructions.
