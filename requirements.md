# Requirements

This document contains the fictional requirements to be fulfilled.

## Description

When you open the CrazyCab app as a passenger, you are able to see a few drivers surrounding you.
These drivers are usually displayed as a car icon. For Halloween, we want to have have some fun, displaying a zombie icon instead of the usual car icon for specific drivers.

CrazyCab drivers send their current coordinates to the backend every five seconds. Our application will use those location updates to differentiate between living and zombie drivers, based on a specific predicate (see below).

To support our growth, we have taken the microservice route. So let’s tackle the basics with a HTTP gateway that either forwards requests or transforms them into [NSQ](https://github.com/nsqio/nsq) messages for asynchronous processing. Then we'll add services that perform tasks related to people transportation.

The task is to implement three services as follows:

* A `gateway` service that either forwards or transforms requests to be respectively processed synchronously or asynchronously.
* A `driver location` service that consumes location update events and stores them.
* A `zombie driver` service that allows users to check whether a driver is a zombie or not.

### 1. Gateway Service

The `Gateway` service is a _public facing service_.
HTTP requests hitting this service are either transformed into [NSQ](https://github.com/nsqio/nsq) messages or forwarded via HTTP to specific services.

The service must be configurable dynamically by loading the provided `gateway/config.yaml` file to register endpoints during its initialization.

#### Public Endpoints

`PATCH /drivers/:id/locations`

**Payload**

```json
{
  "latitude": 48.864193,
  "longitude": 2.350498
}
```

**Role:**

During a typical day, thousands of drivers send their coordinates every 5 seconds to this endpoint.

**Behaviour**

Coordinates received on this endpoint are converted to [NSQ](https://github.com/nsqio/nsq) messages listened by the `Driver Location` service.

---

`GET /drivers/:id`

**Response**

```json
{
  "id": 42,
  "zombie": true
}
```

**Role:**

Users request this endpoint to know if a driver is a zombie.
A driver is a zombie if he has driven less than 500 meters in the last 5 minutes.

**Behaviour**

This endpoint forwards the HTTP request to the `Zombie Driver` service.

### 2. Driver Location Service

The `Driver Location` service is a microservice that consumes drivers' location messages published by the `Gateway` service and stores them in a Redis database.

It also provides an internal endpoint that allows other services to retrieve the drivers' locations, filtered and sorted by their addition date.

#### Internal Endpoint

`GET /drivers/:id/locations?minutes=5`

**Response**

```json
[
  {
    "latitude": 48.864193,
    "longitude": 2.350498,
    "updated_at": "2018-04-05T22:36:16Z"
  },
  {
    "latitude": 48.863921,
    "longitude":  2.349211,
    "updated_at": "2018-04-05T22:36:21Z"
  }
]
```

**Role:**

This endpoint is called by the `Zombie Driver` service.

**Behaviour**

For a given driver, returns all the locations from the last 5 minutes (given `minutes=5`).


### 3. Zombie Driver Service

The `Zombie Driver` service is a microservice that determines if a driver is a zombie or not.

#### Internal Endpoint

`GET /drivers/:id`

**Response**

```
{
  "id": 42,
  "zombie": true
}
```

**Role:**

This endpoint is called by the `Gateway` service.

**Predicate**

> A driver is a zombie if he has driven less than 500 meters in the last 5 minutes.


It would be desirable (bonus point) that the set of rules that determines if a driver is a zombie can be easily changed. For example, a zombie is a driver that hasn't moved more than 1km over the last 15 minutes.


**Behaviour**

Returns the zombie state of a given driver.


## Prerequisites

* Handle all failure cases..
* The gateway should be configured using the `gateway/config.yaml` file
* Provide a clear explanation of the approach and design choices.
* Provide a proper `README.md`:
  * Explaining how to setup and run the code.
  * Including all information which could be considered useful for a seamless on-boarding a new team mate.


NOTE all the resources provided on the initial commit can be used.

## Bonus

* Add metrics / request tracing / circuit breaker.
* Add whatever you think is necessary to make the app amazing.
