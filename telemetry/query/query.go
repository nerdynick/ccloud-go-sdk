package query

import (
	"encoding/json"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/agg"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/filter"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/granularity"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/group"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/interval"
)

const (
	//LifecycleStagePreview is a static def for referencing the PREVIEW Lifecycle Stage of a metric
	LifecycleStagePreview string = "PREVIEW"
	//LifecycleStageGeneral is a static def for referencing the GENERAL_AVAILABILITY Lifecycle Stage of a metric
	LifecycleStageGeneral string = "GENERAL_AVAILABILITY"
)

// Query to Confluent Telemetry API metric endpoint
// This is the JSON structure for the endpoint
// https://api.telemetry.confluent.cloud/v1/metrics/cloud/descriptors
type Query struct {
	Aggreations []agg.Aggregation       `json:"aggregations,omitempty"`
	Filter      filter.Filter           `json:"filter,omitempty"`
	Granularity granularity.Granularity `json:"granularity,omitempty"`
	GroupBy     group.Group             `json:"group_by,omitempty"`
	Intervals   []interval.Interval     `json:"intervals,omitempty"`
	Limit       int                     `json:"limit,omitempty"`
	Metric      metric.Metric           `json:"metric,omitempty"`
}

// func (q *Query) Limit(limit int) *Query {
// 	q.Limit = limit
// }

//ToJSON is a utility function to marshel the Query into a JSON String
func (query Query) ToJSON() ([]byte, error) {
	return json.Marshal(query)
}

func (query Query) Validate() error {
	for _, a := range query.Aggreations {
		err := a.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
