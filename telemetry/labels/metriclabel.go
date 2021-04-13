package labels

import "encoding/json"

var (
	//MetricTopic is a static def for the Topic Label
	MetricTopic Metric = NewMetric("metric.topic")
	//MetricType is a static def for the Type Label
	MetricType Metric = NewMetric("metric.type")
	//MetricPartition is a static def for the Partition Label
	MetricPartition Metric = NewMetric("metric.partition")

	//KnownMetrics is a collection of all the available MetricLabels
	KnownMetrics []Metric = []Metric{
		MetricTopic,
		MetricType,
		MetricPartition,
	}
)

//Metric string type to extend extra helper functions
type Metric struct {
	Key  string `json:"key" json:"key"`
	Desc string `json:"description" json:"description"`
}

func (l Metric) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Key)
}
func (l Metric) String() string {
	return l.Key
}

//NewMetric creates a new Metric Label using a given name
func NewMetric(name string) Metric {
	return Metric{
		Key: name,
	}
}

//IsValid checks in the current label is a valid, available, and known label
func (m Metric) IsValid() bool {
	for _, l := range KnownMetrics {
		if m.Equals(l) {
			return true
		}
	}
	return false
}

//Equals tests if the MetricLabel is equal to a string representation
func (m Metric) Equals(label Metric) bool {
	if m.Key == label.Key {
		return true
	}
	return false
}
