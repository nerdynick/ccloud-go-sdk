package ccloudmetrics

import (
	"time"
)

//AvailableMetricResponse is a struct to house the full response from a GetAvailableMetrics() or GetCurrentlyAvailableMetrics() call
type AvailableMetricResponse struct {
	AvailableMetrics []Metric `json:"data" cjson:"data"`
	Meta             Meta     `json:"meta,omitempty" cjson:"meta,omitempty"`
}

//MetaPagination is a struct to house the Pagination information for a given result
type MetaPagination struct {
	PageSize  int `json:"page_size"`
	TotalSize int `json:"total_size,omitempty"`
}

//Meta is a struct to house the Meta information for a given result
type Meta struct {
	Pagination MetaPagination `json:"pagination"`
}

//QueryData is a struct that represents a given query result's data point
type QueryData struct {
	Timestamp string  `json:"timestamp,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Topic     string  `json:"metric.label.topic,omitempty"`
	Cluster   string  `json:"metric.label.cluster_id,omitempty"`
	Type      string  `json:"metric.label.type,omitempty"`
	Partition string  `json:"metric.label.partition,omitempty"`
}

//Time returned the Timestamp parsed into time.Time
func (q QueryData) Time() (time.Time, error) {
	return time.Parse(TimeFormatStr, q.Timestamp)
}

//HasCluster checks to see if a Cluster Agg was preformed
func (q QueryData) HasCluster() bool {
	if q.Cluster != "" {
		return true
	}
	return false
}

//HasTopic checks to see if a Topic Agg was preformed
func (q QueryData) HasTopic() bool {
	if q.Topic != "" {
		return true
	}
	return false
}

//HasType checks to see if a Type Agg was preformed
func (q QueryData) HasType() bool {
	if q.Type != "" {
		return true
	}
	return false
}

//HasPartition checks to see if a Partition Agg was preformed
func (q QueryData) HasPartition() bool {
	if q.Partition != "" {
		return true
	}
	return false
}

//QueryResponse is a struct that represents a given query result
type QueryResponse struct {
	Data []QueryData `json:"data"`
	Meta Meta        `json:"meta,omitempty"`
}

//ErrorResponse when a none 200 HTTP Status is returned. This handles the JSON
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

//Error a given error
type Error struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}
