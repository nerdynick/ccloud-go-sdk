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

	GranularityOneMin     string = "PT1M"
	GranularityFiveMin    string = "PT5M"
	GranularityFifteenMin string = "PT15M"
	GranularityThirtyMin  string = "PT30M"
	GranularityOneHour    string = "PT1H"
	GranularityAll        string = "ALL"

	LifecycleStagePreview string = "PREVIEW"
	LifecycleStageGeneral string = "GENERAL_AVAILABILITY"
)

var (
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

func NewClusterFilter(cluster string) Filter {
	return Filter{
		Field: MetricLabelCluster,
		Op:    OpEq,
		Value: cluster,
	}
}
func NewTopicFilter(topic string) Filter {
	return Filter{
		Field: MetricLabelTopic,
		Op:    OpEq,
		Value: topic,
	}
}
func NewTypeFilter(ty string) Filter {
	return Filter{
		Field: MetricLabelType,
		Op:    OpEq,
		Value: ty,
	}
}

func NewFilterCollection(op string, filters ...Filter) FilterHeader {
	return FilterHeader{
		Op:      op,
		Filters: filters,
	}
}

func NewTimeInterval(startTime time.Time, endTime time.Time) string {
	return fmt.Sprintf("%s/%s", startTime.Format(TimeFormatStr), endTime.Format(TimeFormatStr))
}

func NewMetricAgg(metric string) Aggregation {
	return Aggregation{
		Agg:    AggSum,
		Metric: metric,
	}
}
