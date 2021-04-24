package resourcetype

import (
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
)

var (
	ResourceTypeKafka          ResourceType = NewResourceType("kafka", metric.KnownKafkaServerMetrics, labels.ResourceKafka)
	ResourceTypeConnector      ResourceType = NewResourceType("connector", metric.KnownConnectorMetrics, labels.ResourceConnector)
	ResourceTypeKSQL           ResourceType = NewResourceType("ksql", metric.KnownKSQLMetrics, labels.ResourceKSQL)
	ResourceTypeSchemaRegistry ResourceType = NewResourceType("schema_registry", metric.KnownSchemaRegMetrics, labels.ResourceSchemaRegistry)
)

//ResourceType represents a returned Resource Type from the API
type ResourceType struct {
	Type         string            `json:"type" cjson:"type"`
	Desc         string            `json:"description,omitempty" cjson:"description,omitempty"`
	Labels       []labels.Resource `json:"labels,omitempty" cjson:"labels,omitempty"`
	KnownMetrics []metric.Metric   `json:"metrics,omitempty"`
}

func NewResourceType(t string, metrics []metric.Metric, labels ...labels.Resource) ResourceType {
	return ResourceType{
		Type:         t,
		Labels:       labels,
		KnownMetrics: metrics,
	}
}
