# confluent-cloud-metrics-go-sdk
A Confluent Cloud GOLang SDK for the Telemetry API


[![Go Report Card](https://goreportcard.com/badge/github.com/nerdynick/confluent-cloud-metrics-go-sdk)](https://goreportcard.com/report/github.com/nerdynick/confluent-cloud-metrics-go-sdk)
[![Build Status](https://travis-ci.com/nerdynick/confluent-cloud-metrics-go-sdk.svg?branch=master)](https://travis-ci.com/nerdynick/confluent-cloud-metrics-go-sdk) 
![GoDoc](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk?status.svg)

# How To

## Install/Get

```bash
go get github.com/nerdynick/confluent-cloud-metrics-go-sdk
```

## Create Telemetry Client

```go
import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry"

func main(){
    telemetryClient := telemetry.New(MyAPIKey, MyAPISecret)
}
```

The `TelemetryClient` also has the following additional attributes that you can adjust to configure how the `TelemetryClient` interacts with the API

```go
PageLimit  int
DataSet    Dataset
BaseURL    string
MaxWorkers int
```

## Get All Available Resources

```go
import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry"

func main(){
    telemetryClient := telemetry.New(MyAPIKey, MyAPISecret)
    resourceTypes, err := telemetryClient.GetAvailableResources()
}
```

## Get All Available Metrics for a Kafka Cluster

```go
import (
    "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry"
    "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry/labels"
)

func main(){
    telemetryClient := telemetry.New(MyAPIKey, MyAPISecret)
    metrics, err := telemetryClient.GetAvailableMetricsForResource(labels.ResourceKafka, "MyClusterID")
}
```

## Get All Available Metrics

_Do Note:_ This func is deprecated post Telemetry API V1. You will want to use the `GetAvailableMetricsForResource` instead.

```go
import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry"

func main(){
    telemetryClient := telemetry.New(MyAPIKey, MyAPISecret)
    metrics, err := telemetryClient.GetAvailableMetrics()
}
```

# Documentation

[Full Docs](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk) | 
[Telemetry Docs](https://godoc.org/github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry) | 