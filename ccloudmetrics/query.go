package ccloudmetrics

import (
	"encoding/json"
	"fmt"
	"strings"
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
	MetricLabelCluster MetricLabel = "metric.label.cluster_id"
	//MetricLabelTopic is a static def for the Topic Label
	MetricLabelTopic MetricLabel = "metric.label.topic"
	//MetricLabelType is a static def for the Type Label
	MetricLabelType MetricLabel = "metric.label.type"
	//MetricLabelPartition is a static def for the Partition Label
	MetricLabelPartition MetricLabel = "metric.label.partition"

	//GranularityOneMin is a static def for the 1 minute grantularity string
	GranularityOneMin Granularity = "PT1M"
	//GranularityFiveMin is a static def for the 5 minute grantularity string
	GranularityFiveMin Granularity = "PT5M"
	//GranularityFifteenMin is a static def for the 15 minute grantularity string
	GranularityFifteenMin Granularity = "PT15M"
	//GranularityThirtyMin is a static def for the 30 minute grantularity string
	GranularityThirtyMin Granularity = "PT30M"
	//GranularityOneHour is a static def for the 1 hour grantularity string
	GranularityOneHour Granularity = "PT1H"
	//GranularityAll is a static def for the ALL grantularity string
	GranularityAll Granularity = "ALL"

	//LifecycleStagePreview is a static def for referencing the PREVIEW Lifecycle Stage of a metric
	LifecycleStagePreview string = "PREVIEW"
	//LifecycleStageGeneral is a static def for referencing the GENERAL_AVAILABILITY Lifecycle Stage of a metric
	LifecycleStageGeneral string = "GENERAL_AVAILABILITY"
)

var (
	//AvailableMetricLabels is a collection of all the available MetricLabels
	AvailableMetricLabels []MetricLabel = []MetricLabel{
		MetricLabelCluster,
		MetricLabelTopic,
		MetricLabelType,
		MetricLabelPartition,
	}
	//AvailableGranularities is a collection of all available Granularities
	AvailableGranularities []Granularity = []Granularity{
		GranularityOneMin,
		GranularityFiveMin,
		GranularityFifteenMin,
		GranularityThirtyMin,
		GranularityOneHour,
		GranularityAll,
	}
)

//MetricLabel string type to extend extra helper functions
type MetricLabel string

//GetFullName returns the full name for a given label
func (m MetricLabel) GetFullName() string {
	return string(m)
}

//GetSimpleName returns a simple name for a given label
func (m MetricLabel) GetSimpleName() string {
	return strings.TrimSuffix(m.GetFullName(), "metric.label.")
}

//IsValid checks in the current label is a valid, available, and known label
func (m MetricLabel) IsValid() bool {
	for _, l := range AvailableMetricLabels {
		if l == m {
			return true
		}
	}
	return false
}

//MetricLabelFromName find a label for a given full name or simple name string
func MetricLabelFromName(name string) MetricLabel {
	if strings.HasPrefix(name, "metric.label") {
		return MetricLabel(name)
	}
	return MetricLabel("metric.label." + name)
}

//Granularity string type to extend extra helper functions
type Granularity string

//IsValid checks in the current Granularity is a valid, available, and known Granularity
func (g Granularity) IsValid() bool {
	for _, l := range AvailableGranularities {
		if l == g {
			return true
		}
	}
	return false
}

//GetStartTimeFromGranularity is a utility func to get a Start time given a granularity
func (g Granularity) GetStartTimeFromGranularity(t time.Time, gran string) time.Time {
	switch g {
	case GranularityOneMin:
		return t.Add(time.Duration(-1) * time.Minute)
	case GranularityFiveMin:
		return t.Add(time.Duration(-5) * time.Minute)
	case GranularityFifteenMin:
		return t.Add(time.Duration(-15) * time.Minute)
	case GranularityThirtyMin:
		return t.Add(time.Duration(-30) * time.Minute)
	case GranularityOneHour:
		return t.Add(time.Duration(-1) * time.Hour)
	}
	return time.Now()
}

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
		Field: MetricLabelCluster.GetFullName(),
		Op:    OpEq,
		Value: cluster,
	}
}

//NewTopicFilter is a utility func to create a Filter for a given topic
func NewTopicFilter(topic string) Filter {
	return Filter{
		Field: MetricLabelTopic.GetFullName(),
		Op:    OpEq,
		Value: topic,
	}
}

//NewTypeFilter is a utility func to create a Filter for a given type
func NewTypeFilter(ty string) Filter {
	return Filter{
		Field: MetricLabelType.GetFullName(),
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
