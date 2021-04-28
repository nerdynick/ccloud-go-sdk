package telemetry

import (
	"encoding/json"
	"errors"

	"github.com/nerdynick/ccloud-go-sdk/logging"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/filter"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/group"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/interval"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/response"
	"go.uber.org/zap"
)

func (client *TelemetryClient) PostLabelQuery(query query.Query) (response.Query, error) {
	url := APIPathAttributes.Format(*client, 2)
	response := response.Query{}

	if len(query.GroupBy.Labels) <= 0 {
		return response, errors.New("Group By is a Required Field for Label Query Types")
	}

	err := client.PostQuery(&response, url, query)
	if err != nil {
		return response, err
	}

	if client.Log.Core().Enabled(logging.InfoLevel) {
		qJson, _ := query.ToJSON()
		resJson, _ := json.Marshal(response)
		client.Log.Info("Query - Response",
			zap.String("URI", url),
			zap.ByteString("Query", qJson),
			zap.ByteString("Response", resJson),
		)
	}

	return response, nil
}

func (client TelemetryClient) LabelQuery(resourceType labels.Resource, resourceID string, metric metric.Metric, field labels.Label, inter interval.Interval) ([]string, error) {
	query := query.Query{
		Filter:    filter.EqualTo(resourceType, resourceID),
		GroupBy:   group.Of(field),
		Intervals: interval.Of(inter),
		Metric:    metric,
	}

	response, err := client.PostLabelQuery(query)
	if err != nil {
		return nil, err
	}

	values := make([]string, len(response.Data))
	for i := 0; i < len(response.Data); i++ {
		values[i] = response.Data[i].Fields[field.String()].(string)
	}

	return values, nil
}

//GetKafkaTopicsForMetric returns all the available topics for a given metric within a window of time
func (client TelemetryClient) GetKafkaTopicsForMetric(cluster string, metric metric.Metric, inter interval.Interval) ([]string, error) {
	return client.LabelQuery(labels.ResourceKafka, cluster, metric, labels.MetricTopic, inter)
}

//GetKafkaRequestTypes returns all the available request types for a given Kafka Cluster
func (client TelemetryClient) GetKafkaRequestTypes(cluster string, inter interval.Interval) ([]string, error) {
	return client.LabelQuery(labels.ResourceKafka, cluster, metric.KafkaServerRequests, labels.MetricType, inter)
}
