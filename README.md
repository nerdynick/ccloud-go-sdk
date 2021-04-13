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