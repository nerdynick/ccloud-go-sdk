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
	KafkaServerReceivedBytes     = New("io.confluent.kafka.server/received_bytes")
	KafkaServerSentBytes         = New("io.confluent.kafka.server/sent_bytes")
	KafkaServerReceivedRecords   = New("io.confluent.kafka.server/received_records")
	KafkaServerSentRecords       = New("io.confluent.kafka.server/sent_records")
	KafkaServerRetainedBytes     = New("io.confluent.kafka.server/retained_bytes")
	KafkaServerActiveConnections = New("io.confluent.kafka.server/active_connection_count")
	KafkaServerRequests          = New("io.confluent.kafka.server/request_count")
	KafkaServerPartition         = New("io.confluent.kafka.server/partition_count")
	KafkaServerSuccessAuth       = New("io.confluent.kafka.server/successful_authentication_count")

	KSQLStreamingUnitCount = New("io.confluent.kafka.ksql/streaming_unit_count")

	SchemaRegSchemaCount = New("io.confluent.kafka.schema_registry/schema_count")

	ConnectorSentRecords            = New("io.confluent.kafka.connect/sent_records")
	ConnectorReceivedRecords        = New("io.confluent.kafka.connect/received_records")
	ConnectorSentBytes              = New("io.confluent.kafka.connect/sent_bytes")
	ConnectorReceivedBytes          = New("io.confluent.kafka.connect/received_bytes")
	ConnectorDeadLetterQueueRecords = New("io.confluent.kafka.connect/dead_letter_queue_records")

	KnownKafkaServerMetrics = []Metric{
		KafkaServerReceivedBytes,
		KafkaServerSentBytes,
		KafkaServerReceivedRecords,
		KafkaServerSentRecords,
		KafkaServerRetainedBytes,
		KafkaServerActiveConnections,
		KafkaServerRequests,
		KafkaServerPartition,
		KafkaServerSuccessAuth,
	}

	KnownKSQLMetrics = []Metric{
		KSQLStreamingUnitCount,
	}

	KnownSchemaRegMetrics = []Metric{
		SchemaRegSchemaCount,
	}

	KnownConnectorMetrics = []Metric{
		ConnectorSentRecords,
		ConnectorReceivedRecords,
		ConnectorSentBytes,
		ConnectorReceivedBytes,
		ConnectorDeadLetterQueueRecords,
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

func (m Metric) MarshalJSON() ([]byte, error) {
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

func New(name string) Metric {
	if strings.HasPrefix(name, metricPrefix) {
		return Metric{Name: name}
	} else {
		return Metric{Name: metricPrefix + name}
	}

}
