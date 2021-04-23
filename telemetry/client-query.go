package telemetry

import (
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/nerdynick/ccloud-go-sdk/logging"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/agg"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/filter"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/granularity"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/group"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/interval"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/response"
	"go.uber.org/zap"
)

func (client TelemetryClient) PostQuery(response interface{}, url string, q query.Query) error {
	if client.Log.Core().Enabled(logging.InfoLevel) {
		qJson, _ := q.ToJSON()
		client.Log.Info("Query - Posting",
			zap.String("URI", url),
			zap.Binary("Query", qJson),
		)
	}

	err := q.Validate()
	if err != nil {
		return err
	}

	return client.Post(&response, url, q)
}

func (client *TelemetryClient) Query(query query.Query) (response.Query, error) {
	url := APIPathDescriptor.Format(*client, 2)
	response := response.Query{}

	err := client.PostQuery(&response, url, query)
	if err != nil {
		return response, err
	}

	if client.Log.Core().Enabled(logging.InfoLevel) {
		qJson, _ := query.ToJSON()
		resJson, _ := json.Marshal(response)
		client.Log.Info("Query - Response",
			zap.String("URI", url),
			zap.Binary("Query", qJson),
			zap.Binary("Response", resJson),
		)
	}

	return response, nil
}

func (client *TelemetryClient) QueryAsync(queryChan <-chan query.Query, resultsChan chan<- response.Query, errsChan chan<- error) {
	for q := range queryChan {
		r, e := client.Query(q)
		resultsChan <- r
		errsChan <- e
	}
}

//QueryMetric returns all the data points for a given metric, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetric(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric) ([]response.Telemetry, error) {
	query := query.Query{
		Filter:      filter.EqualTo(resourceType, resourceID),
		Intervals:   interval.Of(inter),
		Aggreations: agg.Of(agg.SumOf(metric)),
		Granularity: granularity,
		GroupBy:     group.Of(resourceType),
		Limit:       client.PageLimit,
	}

	response, err := client.Query(query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricAsync returns all the data points for a given metric, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetricAsync(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metricChan <-chan metric.Metric, resultsChan chan<- map[string][]response.Telemetry, errsChan chan<- map[string]error) {
	for metric := range metricChan {
		r, e := client.QueryMetric(resourceType, resourceID, granularity, inter, metric)
		if e != nil {
			err := map[string]error{}
			err[metric.Name] = e
			errsChan <- err
		} else {
			res := map[string][]response.Telemetry{}
			res[metric.Name] = r
			resultsChan <- res
		}
	}
}

//QueryMetrics returns all the data points for a given metrics, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetrics(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, timeout time.Duration, metrics ...metric.Metric) (map[string][]response.Telemetry, map[string]error) {
	numMetrics := len(metrics)
	metricsChan := make(chan metric.Metric, numMetrics)
	resultsChan := make(chan map[string][]response.Telemetry, numMetrics)
	errorsChan := make(chan map[string]error, numMetrics)

	defer close(resultsChan)
	defer close(errorsChan)

	client.Log.Debug("Starting up routines")
	for id := 0; id < int(math.Min(float64(numMetrics), float64(client.MaxWorkers))); id++ {
		go client.QueryMetricAsync(resourceType, resourceID, granularity, inter, metricsChan, resultsChan, errorsChan)
	}

	client.Log.Debug("Sending Metrics")
	for _, metric := range metrics {
		client.Log.Debug("Sending Metric: " + metric.Name)
		metricsChan <- metric
	}
	client.Log.Debug("Done Sending Metrics. Closing Channel")
	close(metricsChan)

	results := make(map[string][]response.Telemetry)
	errors := make(map[string]error)
	for {
		select {
		case r := <-resultsChan:
			for m, re := range r {
				results[m] = append(results[m], re...)
			}
		case e := <-errorsChan:
			for m, er := range e {
				errors[m] = er
			}
		case <-time.After(timeout):

		}
	}

	return results, errors
}

//QueryMetricAndLabel returns all the data points for a given metric, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetricAndLabel(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, lbl labels.Metric, lblValue string) ([]response.Telemetry, error) {
	query := query.Query{
		Filter:      filter.EqualTo(resourceType, resourceID).AndEqualTo(lbl, lblValue),
		Intervals:   interval.Of(inter),
		Aggreations: agg.Of(agg.SumOf(metric)),
		Granularity: granularity,
		GroupBy:     group.Of(resourceType).And(lbl),
		Limit:       client.PageLimit,
	}

	response, err := client.Query(query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricAndType returns all the data points for a given metric and type, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetricAndType(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, reqType string) ([]response.Telemetry, error) {
	return client.QueryMetricAndLabel(resourceType, resourceID, granularity, inter, metric, labels.MetricType, reqType)
}

//QueryMetricAndTopic returns all the data points for a given metric and topic, aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetricAndTopic(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, topic string) ([]response.Telemetry, error) {
	if topic == "*" || strings.ToLower(topic) == "all" {
		return client.QueryMetricForAllTopics(resourceType, resourceID, granularity, inter, metric)
	}
	return client.QueryMetricAndLabel(resourceType, resourceID, granularity, inter, metric, labels.MetricTopic, topic)
}

//QueryMetricAndTopicWithPartitions returns all the data points for a given metric and topic, aggregated up to the given granularity, within the given window of time, including aggregations to the partition
func (client *TelemetryClient) QueryMetricAndTopicWithPartitions(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric, topic string) ([]response.Telemetry, error) {
	query := query.Query{
		Filter:      filter.EqualTo(resourceType, resourceID),
		Intervals:   interval.Of(inter),
		Aggreations: agg.Of(agg.SumOf(metric)),
		Granularity: granularity,
		GroupBy:     group.Of(resourceType).And(labels.MetricTopic).And(labels.MetricPartition),
		Limit:       client.PageLimit,
	}

	response, err := client.Query(query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricForAllTopics returns all the data points, fetched in parallel, for a given metric and all available topics (As returned by GetTopicsForMetric), aggregated up to the given granularity, within the given window of time
func (client *TelemetryClient) QueryMetricForAllTopics(resourceType labels.Resource, resourceID string, granularity granularity.Granularity, inter interval.Interval, metric metric.Metric) ([]response.Telemetry, error) {
	query := query.Query{
		Filter:      filter.EqualTo(resourceType, resourceID),
		Intervals:   interval.Of(inter),
		Aggreations: agg.Of(agg.SumOf(metric)),
		Granularity: granularity,
		GroupBy:     group.Of(resourceType).And(labels.MetricTopic),
		Limit:       client.PageLimit,
	}

	response, err := client.Query(query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}
