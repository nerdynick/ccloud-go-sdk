package ccloudmetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	//TimeFormatStr is the format string to use when formating Time for use within Query Intervals
	TimeFormatStr string = time.RFC3339

	//OpAnd is a static def for AND Operand
	OpAnd string = "AND"
	//OpEq is a static def for EQ Operand
	OpEq string = "EQ"
	//OpOr is a static def for OR Operand
	OpOr string = "OR"

	//AggSum is a static def for SUM Aggrigation
	AggSum string = "SUM"

	//LifecycleStagePreview is a static def for referencing the PREVIEW Lifecycle Stage of a metric
	LifecycleStagePreview string = "PREVIEW"
	//LifecycleStageGeneral is a static def for referencing the GENERAL_AVAILABILITY Lifecycle Stage of a metric
	LifecycleStageGeneral string = "GENERAL_AVAILABILITY"
)

// Query to Confluent Cloud API metric endpoint
// This is the JSON structure for the endpoint
// https://api.telemetry.confluent.cloud/v1/metrics/cloud/descriptors
type Query struct {
	Aggreations []Aggregation `json:"aggregations,omitempty"`
	Filter      FilterHeader  `json:"filter,omitempty"`
	Granularity string        `json:"granularity,omitempty"`
	GroupBy     []string      `json:"group_by,omitempty"`
	Intervals   []string      `json:"intervals,omitempty"`
	Limit       int           `json:"limit,omitempty"`
	Metric      string        `json:"metric,omitempty"`
}

//ToJSON is a utility function to marshel the Query into a JSON String
func (query Query) ToJSON() ([]byte, error) {
	return json.Marshal(query)
}

func (query Query) Validate() error {
	for _, a := range query.Aggreations {
		err := a.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

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

// FilterHeader to use for a query
type FilterHeader struct {
	Op      string   `json:"op"`
	Filters []Filter `json:"filters"`
}

// Filter structure
type Filter struct {
	Field   string   `json:"field,omitempty"`
	Op      string   `json:"op"`
	Value   string   `json:"value,omitempty"`
	Filters []Filter `json:"filters,omitempty"`
}

//NewClusterFilter is a utility func to create a Filter for a given cluster
func NewClusterFilter(cluster string) Filter {
	return Filter{
		Field: MetricLabelCluster.String(),
		Op:    OpEq,
		Value: cluster,
	}
}

//NewTopicFilter is a utility func to create a Filter for a given topic
func NewTopicFilter(topic string) Filter {
	return Filter{
		Field: MetricLabelTopic.String(),
		Op:    OpEq,
		Value: topic,
	}
}

//NewTypeFilter is a utility func to create a Filter for a given type
func NewTypeFilter(ty string) Filter {
	return Filter{
		Field: MetricLabelType.String(),
		Op:    OpEq,
		Value: ty,
	}
}

//NewFilterCollection is a utility func to create a Filter for a collection of Filters
func NewFilterCollection(op string, filters ...Filter) FilterHeader {
	return FilterHeader{
		Op:      op,
		Filters: filters,
	}
}

//NewTimeInterval is a utility func to create a Interval String for a given Time range
func NewTimeInterval(startTime time.Time, endTime time.Time) string {
	return fmt.Sprintf("%s/%s", startTime.Format(TimeFormatStr), endTime.Format(TimeFormatStr))
}

//NewMetricAgg is a utility func to create a Aggregation for a given metric
func NewMetricAgg(metric string) Aggregation {
	return Aggregation{
		Agg:    AggSum,
		Metric: metric,
	}
}
