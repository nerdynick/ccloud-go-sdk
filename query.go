package ccloudmetrics

import "encoding/json"

// Query to Confluent Cloud API metric endpoint
// This is the JSON structure for the endpoint
// https://api.telemetry.confluent.cloud/v1/metrics/cloud/descriptors
type Query struct {
	Aggreations []Aggregation `json:"aggregations,omitempty"`
	Filter      FilterHeader  `json:"filter,omitempty"`
	Granularity string        `json:"granularity,omitempty"`
	GroupBy     []string      `json:"group_by,omitempty"`
	Intervals   []string      `json:"intervals,omitempty"`
	Limit       int           `json:"limit,omitempty"`
	Metric      string        `json:"metric,omitempty"`
}

func (query Query) ToJSON() ([]byte, error) {
	return json.Marshal(query)
}

// Aggregation for a Confluent Cloud API metric
type Aggregation struct {
	Agg    string `json:"agg"`
	Metric string `json:"metric"`
}

// FilterHeader to use for a query
type FilterHeader struct {
	Op      string   `json:"op"`
	Filters []Filter `json:"filters"`
}

// Filter structure
type Filter struct {
	Field   string   `json:"field,omitempty"`
	Op      string   `json:"op"`
	Value   string   `json:"value,omitempty"`
	Filters []Filter `json:"filters,omitempty"`
}
