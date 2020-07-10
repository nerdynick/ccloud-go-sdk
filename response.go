package ccloudmetrics

type AvailableMetricLabel struct {
	Name *string `json:"key"`
	Desc *string `json:"description"`
}
type AvailableMetric struct {
	Name           *string                 `json:"name"`
	Desc           *string                 `json:"description"`
	Type           *string                 `json:"type"`
	LifecycleStage *string                 `json:"lifecycle_stage"`
	Labels         *[]AvailableMetricLabel `json:"labels"`
}
type AvailableMetricResponse struct {
	AvailableMetrics *[]AvailableMetric `json:"data"`
	Meta             *Meta              `json:"meta,omitempty"`
}

type CurrentlyAvailableMetric struct {
	Name *string `json:"metric"`
}
type CurrentlyAvailableMetricResponse struct {
	AvailableMetrics *[]CurrentlyAvailableMetric `json:"data"`
	Meta             *Meta                       `json:"meta,omitempty"`
}

type MetaPagination struct {
	PageSize  int `json:"page_size"`
	TotalSize int `json:"total_size,omitempty"`
}
type Meta struct {
	Pagination MetaPagination `json:"pagination"`
}

type QueryData struct {
	Timestamp string  `json:"timestamp,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Topic     string  `json:"metric.label.topic,omitempty"`
	Cluster   string  `json:"metric.label.cluster_id,omitempty"`
	Type      string  `json:"metric.label.type,omitempty"`
}
type QueryResponse struct {
	Data []QueryData `json:"data"`
	Meta Meta        `json:"meta,omitempty"`
}
