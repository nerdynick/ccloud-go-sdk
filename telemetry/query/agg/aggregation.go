package agg

import (
	"errors"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry/metric"
)

const (
	//AggSum is a static def for SUM Aggrigation
	AggSum string = "SUM"
)

// Aggregation for a Confluent Cloud API metric
type Aggregation struct {
	Agg    string `json:"agg"`
	Metric string `json:"metric"`
}

func (a Aggregation) Validate() error {
	if a.Agg == "" {
		return errors.New("Agg can not be empty/nil")
	}

	if a.Metric == "" {
		return errors.New("Metric can not be empty/nil")
	}
	return nil
}

func SumOf(metric metric.Metric) Aggregation {
	return Aggregation{
		Agg:    AggSum,
		Metric: metric.Name,
	}
}

func Of(aggs ...Aggregation) []Aggregation {
	return aggs
}
