package ccloudmetrics

import "strings"

const (
	//MetricLabelCluster is a static def for the Cluster ID Label
	MetricLabelCluster MetricLabel = "metric.label.cluster_id"
	//MetricLabelTopic is a static def for the Topic Label
	MetricLabelTopic MetricLabel = "metric.label.topic"
	//MetricLabelType is a static def for the Type Label
	MetricLabelType MetricLabel = "metric.label.type"
	//MetricLabelPartition is a static def for the Partition Label
	MetricLabelPartition MetricLabel = "metric.label.partition"
)

var (
	//AvailableMetricLabels is a collection of all the available MetricLabels
	AvailableMetricLabels []string = []string{
		MetricLabelCluster.String(),
		MetricLabelTopic.String(),
		MetricLabelType.String(),
		MetricLabelPartition.String(),
	}
)

//MetricLabel string type to extend extra helper functions
type MetricLabel string

//IsValid checks in the current label is a valid, available, and known label
func (m MetricLabel) IsValid() bool {
	for _, l := range AvailableMetricLabels {
		if m.Equals(l) {
			return true
		}
	}
	return false
}

//Equals tests if the MetricLabel is equal to a string representation
func (m MetricLabel) Equals(str string) bool {
	if m == NewMetricLabel(str) {
		return true
	}
	return false
}

func (m MetricLabel) String() string {
	return string(m)
}

//ExtendedMetricLabel transforms MetricLabel into an ExtendedMetricLabel
func (m MetricLabel) ExtendedMetricLabel() ExtendedMetricLabel {
	return ExtendedMetricLabel{
		Name: m.String(),
	}
}

//NewMetricLabel creats a new MetricLabel from a string value
func NewMetricLabel(name string) MetricLabel {
	if strings.HasPrefix(name, "metric.label.") {
		return MetricLabel(name)
	}
	return MetricLabel("metric.label." + name)
}
