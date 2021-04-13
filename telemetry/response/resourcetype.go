package response

import "github.com/nerdynick/ccloud-go-sdk/telemetry/labels"

//ResourceType represents a returned Resource Type from the API
type ResourceType struct {
	Type   string            `json:"type" cjson:"type"`
	Desc   string            `json:"description,omitempty" cjson:"description,omitempty"`
	Labels []labels.Resource `json:"labels,omitempty" cjson:"labels,omitempty"`
}
