package ccloudmetrics

import (
	"encoding/json"
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

	//MetricLabelCluster is a static def for the Cluster ID Label
	MetricLabelCluster string = "metric.label.cluster_id"
	//MetricLabelTopic is a static def for the Topic Label
	MetricLabelTopic string = "metric.label.topic"
	//MetricLabelType is a static def for the Type Label
	MetricLabelType string = "metric.label.type"
	//MetricLabelPartition is a static def for the Partition Label
	MetricLabelPartition string = "metric.label.partition"

	//GranularityOneMin is a static def for the 1 minute grantularity string
	GranularityOneMin string = "PT1M"
	//GranularityFiveMin is a static def for the 5 minute grantularity string
	GranularityFiveMin string = "PT5M"
	//GranularityFifteenMin is a static def for the 15 minute grantularity string
	GranularityFifteenMin string = "PT15M"
	//GranularityThirtyMin is a static def for the 30 minute grantularity string
	GranularityThirtyMin string = "PT30M"
	//GranularityOneHour is a static def for the 1 hour grantularity string
	GranularityOneHour string = "PT1H"
	//GranularityAll is a static def for the ALL grantularity string
	GranularityAll string = "ALL"

	//LifecycleStagePreview is a static def for referencing the PREVIEW Lifecycle Stage of a metric
	LifecycleStagePreview string = "PREVIEW"
	//LifecycleStageGeneral is a static def for referencing the GENERAL_AVAILABILITY Lifecycle Stage of a metric
	LifecycleStageGeneral string = "GENERAL_AVAILABILITY"
)

var (
	//AvailableGranularities is a collection of all available Granularities
	AvailableGranularities []string = []string{
		GranularityOneMin,
		GranularityFiveMin,
		GranularityFifteenMin,
		GranularityThirtyMin,
		GranularityOneHour,
		GranularityAll,
	}
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

// Aggregation for a Confluent Cloud API metric
type Aggregation struct {
	Agg    string `json:"agg"`
	Metric string `json:"metric"`
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
		Field: MetricLabelCluster,
		Op:    OpEq,
		Value: cluster,
	}
}

//NewTopicFilter is a utility func to create a Filter for a given topic
func NewTopicFilter(topic string) Filter {
	return Filter{
		Field: MetricLabelTopic,
		Op:    OpEq,
		Value: topic,
	}
}

//NewTypeFilter is a utility func to create a Filter for a given type
func NewTypeFilter(ty string) Filter {
	return Filter{
		Field: MetricLabelType,
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
