package ccloudmetrics

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	metricPrefix = "io.confluent.kafka.server/"
)

var (
	MetricReceivedBytes     = Metric{Name: "io.confluent.kafka.server/received_bytes"}
	MetricSentBytes         = Metric{Name: "io.confluent.kafka.server/sent_bytes"}
	MetricReceivedRecords   = Metric{Name: "io.confluent.kafka.server/received_records"}
	MetricSentRecords       = Metric{Name: "io.confluent.kafka.server/sent_records"}
	MetricRetainedBytes     = Metric{Name: "io.confluent.kafka.server/retained_bytes"}
	MetricActiveConnections = Metric{Name: "io.confluent.kafka.server/active_connection_count"}
	MetricRequests          = Metric{Name: "io.confluent.kafka.server/request_count"}
	MetricPartition         = Metric{Name: "io.confluent.kafka.server/partition_count"}
	MetricSuccessAuth       = Metric{Name: "io.confluent.kafka.server/successful_authentication_count"}
)

//Metric is a struct to house the Metric details for a returned metric
type Metric struct {
	Name           string                `json:"name" cjson:"name"`
	Desc           string                `json:"description,omitempty" cjson:"description,omitempty"`
	Type           string                `json:"type,omitempty" cjson:"type,omitempty"`
	LifecycleStage string                `json:"lifecycle_stage,omitempty" cjson:"lifecycle_stage,omitempty"`
	Labels         []ExtendedMetricLabel `json:"labels,omitempty" cjson:"labels,omitempty"`
}

//Check is a given metric name is equal to this metric
func (m Metric) Matches(name string) bool {
	return m.Name == name || m.ShortName() == name
}

//ShortName returned a simple shorter name, without all the namespacing
func (m Metric) ShortName() string {
	return strings.TrimPrefix(m.Name, metricPrefix)
}

//HasLabel checks if a given AvailableMetric has a given label
func (m Metric) HasLabel(label MetricLabel) bool {
	if m.Labels != nil {
		for _, l := range m.Labels {
			if label.Equals(l.Name) {
				return true
			}
		}
	}
	return false
}

//GetValidLabels given a whitelist of possible labels will return a collection of labels that are valid to use against this metric
func (m Metric) GetValidLabels(whitelist []MetricLabel) []string {
	log.WithFields(log.Fields{
		"AvailableLabels":   m.Labels,
		"WhitelistedLabels": whitelist,
	}).Debug("Getting Valid Labels")

	labels := []string{}
	for _, l := range whitelist {
		if m.HasLabel(l) {
			labels = append(labels, l.String())
		}
	}
	return labels
}

//NewMetric Create a minimilistic Metric representation for use in Query Calls
func NewMetric(name string) Metric {
	if strings.HasPrefix(name, metricPrefix) {
		return Metric{Name: name}
	} else {
		return Metric{Name: metricPrefix + name}
	}

}

//ExtendedMetricLabel is a struct to house the Label details for a return metric
type ExtendedMetricLabel struct {
	Name string `json:"key" json:"key"`
	Desc string `json:"description" json:"description"`
}

//MetricLabel Transform ExtendedMetricLabel into MetricLabel
func (e ExtendedMetricLabel) MetricLabel() MetricLabel {
	return NewMetricLabel(e.Name)
}
