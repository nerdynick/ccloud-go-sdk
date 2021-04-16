package response

import "github.com/nerdynick/ccloud-go-sdk/telemetry/metric"

//BaseResponse is a common struct of fields for all API responses
type BaseResponse struct {
	Meta  Meta  `json:"meta,omitempty"`
	Links Links `json:"links,omitempty"`
}

//Meta is a struct to house the Meta information for a given result
type Meta struct {
	Pagination MetaPagination `json:"pagination"`
}

//MetaPagination is a struct to house the Pagination information for a given result
type MetaPagination struct {
	PageSize  int `json:"page_size"`
	TotalSize int `json:"total_size,omitempty"`
}

//Links represents the Links return data
type Links struct {
	Next string `json:"next,omitempty"`
}

//Query represents a collection of Telemetry records as returned from a Telemetry Query request
type Query struct {
	*BaseResponse
	Data []Telemetry `json:"data"`
}

//Resources respresents a collection of resources as returned from lookup of resources
type Resources struct {
	*BaseResponse
	ResourceTypes []ResourceType `json:"data"`
}

//Metrics represents a collection of possible Available Metrics
type Metrics struct {
	*BaseResponse
	AvailableMetrics []metric.Metric `json:"data"`
}
