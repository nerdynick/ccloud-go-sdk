package labels

import "encoding/json"

var (
	//ResourceKafka is a static reference for a Kafka Cluster's ID
	ResourceKafka Resource = newResource("kafka.id")

	//KnownResources is a collection of known resource labels at this time
	KnownResources []Resource = []Resource{
		ResourceKafka,
	}
)

//Resource struct to represent a Resource Label
type Resource struct {
	Key  string `json:"key" json:"key"`
	Desc string `json:"description" json:"description"`
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
