package group

import (
	"encoding/json"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
)

type Group struct {
	Labels []labels.Label
}

func (g Group) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.Labels)
}
func (g Group) And(labels ...labels.Label) Group {
	g.Labels = append(g.Labels, labels...)
	return g
}

func By(labels ...labels.Label) Group {
	return Of(labels...)
}
func Of(labels ...labels.Label) Group {
	return Group{
		Labels: labels,
	}
}
