package response

import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry/metric"

type baseResponse struct {
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

//Links respresents the Links return data
type Links struct {
	Next string `json:"next,omitempty"`
}

type Query struct {
	*baseResponse
	Data []Telemetry `json:"data"`
}

type Resources struct {
	*baseResponse
	ResourceTypes []ResourceType `json:"data"`
}

type Metrics struct {
	*baseResponse
	AvailableMetrics []metric.Metric `json:"data"`
}
