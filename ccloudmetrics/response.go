package ccloudmetrics

//AvailableMetricLabel is a struct to house the Label details for a return metric
type AvailableMetricLabel struct {
	Name string `json:"key" json:"key"`
	Desc string `json:"description" json:"description"`
}

//AvailableMetric is a struct to house the Metric details for a returned metric
type AvailableMetric struct {
	Name           string                 `json:"name" cjson:"name"`
	Desc           string                 `json:"description,omitempty" cjson:"description,omitempty"`
	Type           string                 `json:"type,omitempty" cjson:"type,omitempty"`
	LifecycleStage string                 `json:"lifecycle_stage,omitempty" cjson:"lifecycle_stage,omitempty"`
	Labels         []AvailableMetricLabel `json:"labels,omitempty" cjson:"labels,omitempty"`
}

//AvailableMetricResponse is a struct to house the full response from a GetAvailableMetrics() or GetCurrentlyAvailableMetrics() call
type AvailableMetricResponse struct {
	AvailableMetrics []AvailableMetric `json:"data" cjson:"data"`
	Meta             Meta              `json:"meta,omitempty" cjson:"meta,omitempty"`
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

//QueryResponse is a struct that represents a given query result
type QueryResponse struct {
	Data []QueryData `json:"data"`
	Meta Meta        `json:"meta,omitempty"`
}
