package telemetry

import (
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/agg"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/filter"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/granularity"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/group"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/interval"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/response"
)

//QueryMetricAndType returns all the data points for a given metric and type, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryKafkaMetricAndType(resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, reqType string) ([]response.Telemetry, error) {
	return client.QueryMetricAndLabel(labels.ResourceKafka, resourceID, granularity, inter, metric, labels.MetricType, reqType)
}

//QueryMetricAndTopic returns all the data points for a given metric and topic, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryKafkaMetricAndTopic(resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, topic string) ([]response.Telemetry, error) {
	if topic == "*" || strings.ToLower(topic) == "all" {
		return client.QueryKafkaMetricForAllTopics(resourceID, granularity, inter, metric)
	}
	return client.QueryMetricAndLabel(labels.ResourceKafka, resourceID, granularity, inter, metric, labels.MetricTopic, topic)
}

//QueryMetricAndTopicWithPartitions returns all the data points for a given metric and topic, aggregated up to the given granularity, within the given window of time, including aggregations to the partition
func (client *TelemetryClient) QueryKafkaMetricAndTopicWithPartitions(resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, topic string) ([]response.Telemetry, error) {
	query := query.Query{
		Filter:       filter.EqualTo(labels.ResourceKafka, resourceID),
		Intervals:    interval.Of(inter),
		Aggregations: agg.Of(agg.SumOf(metric)),
		Granularity:  granularity,
		GroupBy:      group.Of(labels.ResourceKafka).And(labels.MetricTopic).And(labels.MetricPartition),
		Limit:        client.PageLimit,
	}

	response, err := client.PostMetricsQuery(query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricForAllTopics returns all the data points, fetched in parallel, for a given metric and all available topics (As returned by GetTopicsForMetric), aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryKafkaMetricForAllTopics(resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric) ([]response.Telemetry, error) {
	query := query.Query{
		Filter:       filter.EqualTo(labels.ResourceKafka, resourceID),
		Intervals:    interval.Of(inter),
		Aggregations: agg.Of(agg.SumOf(metric)),
		Granularity:  granularity,
		GroupBy:      group.Of(labels.ResourceKafka).And(labels.MetricTopic),
		Limit:        client.PageLimit,
	}

	response, err := client.PostMetricsQuery(query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}
