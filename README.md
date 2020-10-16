# confluent-cloud-metrics-go-sdk
Confluent Cloud Metrics API SDK for GO as well as a CLI for interaction with the Metrics API


[![Go Report Card](https://goreportcard.com/badge/github.com/nerdynick/confluent-cloud-metrics-go-sdk)](https://goreportcard.com/report/github.com/nerdynick/confluent-cloud-metrics-go-sdk)
[![Build Status](https://travis-ci.org/nerdynick/confluent-cloud-metrics-go-sdk.svg?branch=master)](https://travis-ci.org/nerdynick/confluent-cloud-metrics-go-sdk) 
![GoDoc](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk?status.svg)



## SDK

To use the SDK you simply just need to import the SDK into your Go project and create a new instance of the MetricsClient.
The client supporta a number of settings to adjust things to your liking. 
See [ccloud.go](ccloudmetrics/ccloud.go) for further details of all available settings.

```go
import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"

func main(){
    client := ccloudmetrics.NewClientMinimal(MyAPIKey,MyAPISecret)
}
```

### Examples

This is just a simple set of example. 
You can see further usage examples with the included CLI which is located within the [cmd](cmd) folder.

**Get All Available Metrics**

```go
import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"

func main(){
    client := ccloudmetrics.NewClientMinimal(MyAPIKey,MyAPISecret)
    client.GetAvailableMetrics()
    client.GetTopicsForMetric(MyClusterID, "io.confluent.kafka.server/retained_bytes", StartTime, EndTime)
}
```

**Get All Available Topics for a given Metric**

```go
import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"

func main(){
    client := ccloudmetrics.NewClientMinimal(MyAPIKey,MyAPISecret)
    client.GetTopicsForMetric(MyClusterID, "io.confluent.kafka.server/retained_bytes", StartTime, EndTime)
}
```

### Documentation

[Full Docs](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk) | 
[SDK Docs](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk/cmd) | 
[CLI Docs](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk/cmd)

## CLI

Included in this codebase is a CLI that leverages the packaged SDK. 
This is to both provide a CLI experience with the Metrics API, but also provide a source of example usages.

### Installing the CLI

```shell
go get github.com/nerdynick/confluent-cloud-metrics-go-sdk
```

### Usage

The CLI is coded to be interacted with similar to GIT.
It leverages the project [Cobra](https://github.com/spf13/cobra) to provide this experience.

Command Tree

* root
  * list
    * metrics - Query Available Metrics
    * topics - Query Available Topics for a Metric
  * query
    * metric - Query Data for a metric
    * metrics - Query Data for multipule metrics
    * topic - Query Data for a metric and topic
    * topics - Query Data for a metric and a list of topics
      * all - Query Data for a metric and all available topics (Queries run in parallel for each topic)

**List Available Metrics**

```shell
./confluent-cloud-metrics-go-sdk list metrics --apikey MY-KEY --apisecret MY-SECRET
```

**List Available Topics for a given Metric**

```shell
./confluent-cloud-metrics-go-sdk list topics --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID  io.confluent.kafka.server/retained_bytes
```

**Query Cluster level results for a given Metric**

```shell
./confluent-cloud-metrics-go-sdk query metric --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID  io.confluent.kafka.server/retained_bytes
```

**Query Topic level results for a given Metric & Topic**

```shell
./confluent-cloud-metrics-go-sdk query topic --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID --metric io.confluent.kafka.server/retained_bytes  MY-TOPIC
```

**Query Topic Parition level results for a given Metric & Topic**

```shell
./confluent-cloud-metrics-go-sdk query topic --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID --metric io.confluent.kafka.server/retained_bytes MY-TOPIC --partitions
```