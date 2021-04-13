package metric

import (
	"encoding/json"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
)

const (
	metricPrefix = "io.confluent.kafka.server/"
)

var (
	ReceivedBytes     = new("io.confluent.kafka.server/received_bytes")
	SentBytes         = new("io.confluent.kafka.server/sent_bytes")
	ReceivedRecords   = new("io.confluent.kafka.server/received_records")
	SentRecords       = new("io.confluent.kafka.server/sent_records")
	RetainedBytes     = new("io.confluent.kafka.server/retained_bytes")
	ActiveConnections = new("io.confluent.kafka.server/active_connection_count")
	Requests          = new("io.confluent.kafka.server/request_count")
	Partition         = new("io.confluent.kafka.server/partition_count")
	SuccessAuth       = new("io.confluent.kafka.server/successful_authentication_count")

	KnownMetrics = []Metric{
		ReceivedBytes,
		SentBytes,
		ReceivedRecords,
		SentRecords,
		RetainedBytes,
		ActiveConnections,
		Requests,
		Partition,
		SuccessAuth,
	}
)

//Metric is a struct to house the Metric details for a returned metric
type Metric struct {
	Name           string          `json:"name" cjson:"name"`
	Desc           string          `json:"description,omitempty" cjson:"description,omitempty"`
	Type           string          `json:"type,omitempty" cjson:"type,omitempty"`
	LifecycleStage string          `json:"lifecycle_stage,omitempty" cjson:"lifecycle_stage,omitempty"`
	Labels         []labels.Metric `json:"labels,omitempty" cjson:"labels,omitempty"`
}

func (m *Metric) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Name)
}

//Matches check if a given metric name is equal to this metric
func (m Metric) Matches(name string) bool {
	return m.Name == name || m.ShortName() == name
}

//ShortName returned a simple shorter name, without all the namespacing
func (m Metric) ShortName() string {
	return strings.TrimPrefix(m.Name, metricPrefix)
}

func new(name string) Metric {
	if strings.HasPrefix(name, metricPrefix) {
		return Metric{Name: name}
	} else {
		return Metric{Name: metricPrefix + name}
	}

}
