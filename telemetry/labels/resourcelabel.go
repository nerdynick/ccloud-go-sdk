package labels

import "encoding/json"

var (
	//ResourceKafka is a static reference for a Kafka Cluster's ID
	ResourceKafka          Resource = newResource("kafka.id")
	ResourceConnector      Resource = newResource("connector.id")
	ResourceKSQL           Resource = newResource("ksql.id")
	ResourceSchemaRegistry Resource = newResource("schema_registry.id")

	//KnownResources is a collection of known resource labels at this time
	KnownResources []Resource = []Resource{
		ResourceKafka,
		ResourceConnector,
		ResourceKSQL,
		ResourceSchemaRegistry,
	}
)

//Resource struct to represent a Resource Label
type Resource struct {
	Key  string `json:"key"`
	Desc string `json:"description"`
}

func (l Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Key)
}
func (l Resource) String() string {
	return l.Key
}

func newResource(key string) Resource {
	return Resource{
		Key: key,
	}
}
