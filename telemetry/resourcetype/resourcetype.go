package resourcetype

import "github.com/nerdynick/ccloud-go-sdk/telemetry/labels"

var (
	ResourceTypeKafka          ResourceType = newResourceType("kafka", labels.ResourceKafka)
	ResourceTypeConnector      ResourceType = newResourceType("connector", labels.ResourceConnector)
	ResourceTypeKSQL           ResourceType = newResourceType("ksql", labels.ResourceKSQL)
	ResourceTypeSchemaRegistry ResourceType = newResourceType("schema_registry", labels.ResourceSchemaRegistry)
)

//ResourceType represents a returned Resource Type from the API
type ResourceType struct {
	Type   string            `json:"type" cjson:"type"`
	Desc   string            `json:"description,omitempty" cjson:"description,omitempty"`
	Labels []labels.Resource `json:"labels,omitempty" cjson:"labels,omitempty"`
}

func newResourceType(t string, labels ...labels.Resource) ResourceType {
	return ResourceType{
		Type:   t,
		Labels: labels,
	}
}
