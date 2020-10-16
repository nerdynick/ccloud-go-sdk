package ccloudmetrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

const (
	descriptorPath string = "/v1/metrics/cloud/descriptors"
	attributesPath string = "/v1/metrics/cloud/attributes"
	availablePath  string = "/v1/metrics/cloud/available"
	queryPath      string = "/v1/metrics/cloud/query"

	//DefaultQueryLimit the default query limit for results
	DefaultQueryLimit int = 1000
	//DefaultBaseURL is the default Domain and Protocol for quering against the Metrics API
	DefaultBaseURL string = "https://api.telemetry.confluent.cloud"
	//DefaultRequestTimeout is the default number of seconds to wait before considering a Metrics API query/request as timedout
	DefaultRequestTimeout int = 60
	//DefaultUserAgent is the default user agent to send
	DefaultUserAgent string = "ccloud-metrics-sdk/go"

	//DefaultMaxWorkers controls the max number of workers in a given Worker Pool that will be spawned
	DefaultMaxWorkers int = 5
)

var (
	cJSON = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "cjson",
	}.Froze()
)

//APIContext is the Contextual set of configs for all Metrics API calls
type APIContext struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	MaxWorkers int
}

//NewAPIContext creates a new instance of the APIContext loaded with the defaults where possible
func NewAPIContext(apiKey string, apiSecret string) APIContext {
	return APIContext{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    DefaultBaseURL,
		MaxWorkers: DefaultMaxWorkers,
	}
}

//HTTPContext is the Contextual set of configs for the HTTP Client making the calls to the Metrics API
type HTTPContext struct {
	RequestTimeout      int
	UserAgent           string
	HTTPHeaders         map[string]string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

//NewHTTPContext creates a new instance of the HTTPContext loaded with the defaults where possible
func NewHTTPContext() HTTPContext {
	return HTTPContext{
		RequestTimeout:      DefaultRequestTimeout,
		UserAgent:           DefaultUserAgent,
		HTTPHeaders:         nil,
		MaxIdleConns:        DefaultMaxWorkers,
		MaxIdleConnsPerHost: DefaultMaxWorkers,
	}
}

//MetricsClient is the SDK Client for making REST calls to the Confluent Metrics API
type MetricsClient struct {
	apiContext  APIContext
	httpContext HTTPContext
	httpClient  http.Client
}

//NewClientFromContext Used to create a new MetricsClient from the given contexts
func NewClientFromContext(context APIContext, httpContext HTTPContext) MetricsClient {
	log.WithFields(log.Fields{
		"APIContext":  context,
		"HTTPContext": httpContext,
	}).Trace("Creating new Metrics Client")

	httpClient := http.Client{
		Timeout: time.Second * time.Duration(httpContext.RequestTimeout),
		Transport: &http.Transport{
			MaxIdleConns:        httpContext.MaxIdleConns,
			MaxIdleConnsPerHost: httpContext.MaxIdleConnsPerHost,
		},
	}
	client := MetricsClient{
		httpClient:  httpClient,
		apiContext:  context,
		httpContext: httpContext,
	}

	return client
}

//NewClientMinimal Used to create a new MetricsClient from the given minimal set of properties
func NewClientMinimal(apiKey string, apiSecret string) MetricsClient {
	return NewClientFromContext(NewAPIContext(apiKey, apiSecret), NewHTTPContext())
}

//NewClientWithDefaults Used to create a new MetricsClient from the minimal set of properties and using defaults where appropriate
func NewClientWithDefaults(apiKey string, apiSecret string, extraHeaders map[string]string) MetricsClient {
	httpContext := NewHTTPContext()
	httpContext.HTTPHeaders = extraHeaders
	return NewClientFromContext(NewAPIContext(apiKey, apiSecret), httpContext)
}

//NewClient Used to create a new MetricsClient from the full set of properties
func NewClient(apiKey string, apiSecret string, baseURL string, requestTime int, userAgent string, maxWorkers int, extraHeaders map[string]string) MetricsClient {
	return NewClientFromContext(APIContext{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    strOrDefault(baseURL, DefaultBaseURL),
		MaxWorkers: intOrDefault(maxWorkers, DefaultMaxWorkers),
	}, HTTPContext{
		RequestTimeout:      intOrDefault(requestTime, DefaultRequestTimeout),
		UserAgent:           strOrDefault(userAgent, DefaultUserAgent),
		HTTPHeaders:         extraHeaders,
		MaxIdleConnsPerHost: intOrDefault(maxWorkers, DefaultMaxWorkers),
		MaxIdleConns:        intOrDefault(maxWorkers, DefaultMaxWorkers),
	})
}

func strOrDefault(val string, def string) string {
	if val != "" {
		return val
	}
	return def
}
func intOrDefault(val int, def int) int {
	if val > 0 {
		return val
	}
	return def
}

//GetAvailableMetrics returns a collection of all the available metrics and their supported labels among other important meta data
func (client MetricsClient) GetAvailableMetrics() ([]Metric, error) {
	result, err := client.SendGet(descriptorPath)

	if err != nil {
		return nil, err
	}

	response := AvailableMetricResponse{}
	json.Unmarshal(result, &response)

	return response.AvailableMetrics, nil
}

//GetCurrentlyMetrics returns all the currently available metrics and their supported labels among other important meta data
func (client MetricsClient) GetCurrentlyMetrics(cluster string) ([]Metric, error) {
	query := Query{
		Filter: NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
	}
	result, err := client.SendPost(availablePath, query)

	if err != nil {
		return nil, err
	}

	response := AvailableMetricResponse{}
	cJSON.Unmarshal(result, &response)

	return response.AvailableMetrics, nil
}

//GetTopicsForMetric returns all the available topics for a given metric within a window of time
func (client MetricsClient) GetTopicsForMetric(cluster string, metric Metric, startTime time.Time, endTime time.Time) ([]string, error) {
	query := Query{
		Filter:    NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
		GroupBy:   metric.GetValidLabels([]MetricLabel{MetricLabelTopic}),
		Intervals: []string{NewTimeInterval(startTime, endTime)},
		Metric:    metric.Name,
	}
	response, err := client.SendQuery(attributesPath, query)
	if err != nil {
		return nil, err
	}

	topics := make([]string, len(response.Data))
	for i := 0; i < len(response.Data); i++ {
		topics[i] = response.Data[i].Topic
	}

	return topics, nil
}

//QueryMetric returns all the data points for a given metric, aggregated up to the given granularity, within the given window of time
func (client MetricsClient) QueryMetric(cluster string, metric Metric, granularity Granularity, startTime time.Time, endTime time.Time) ([]QueryData, error) {
	query := Query{
		Filter:      NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
		Intervals:   []string{NewTimeInterval(startTime, endTime)},
		Aggreations: []Aggregation{NewMetricAgg(metric.Name)},
		Granularity: granularity.String(),
		GroupBy:     metric.GetValidLabels([]MetricLabel{MetricLabelCluster}),
		Limit:       DefaultQueryLimit,
	}

	response, err := client.SendQuery(queryPath, query)
	log.WithField("response", response).Debug("QueryMetric Response")
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetrics returns all the data points for a given metrics, aggregated up to the given granularity, within the given window of time
func (client MetricsClient) QueryMetrics(cluster string, metric []Metric, granularity Granularity, startTime time.Time, endTime time.Time) ([]QueryData, error) {
	metricsChan := make(chan Metric, len(metric))
	resultsChan := make(chan []QueryData, len(metric))
	errorsChan := make(chan error, len(metric))
	waitGroup := sync.WaitGroup{}

	log.Debug("Starting up routines")
	for id := 0; id < client.apiContext.MaxWorkers; id++ {
		go client.queryMetricWorker(cluster, granularity, startTime, endTime, id, &waitGroup, metricsChan, resultsChan, errorsChan)
	}

	log.Debug("Sending Topics")
	for _, metric := range metric {
		log.Debug("Sending Topic: " + metric.Name)
		metricsChan <- metric
		waitGroup.Add(1)
	}
	log.Debug("Done Sending Metrics. Closing Channel")
	close(metricsChan)

	waitGroup.Wait()
	close(resultsChan)
	close(errorsChan)

	log.Debug("Processing Errors")
	finalErrors := []string{}
	var err error
	for err := range errorsChan {
		finalErrors = append(finalErrors, err.Error())
	}
	if len(finalErrors) > 0 {
		err = errors.New(strings.Join(finalErrors, "\n\n"))
	} else {
		err = nil
	}

	log.Debug("Processing Results")
	queryData := []QueryData{}
	for res := range resultsChan {
		queryData = append(queryData, res...)
	}
	log.WithField("response", queryData).Debug("QueryMetrics Response")
	return queryData, err
}

func (client MetricsClient) queryMetricWorker(cluster string, granularity Granularity, startTime time.Time, endTime time.Time, workerID int, wg *sync.WaitGroup, metrics <-chan Metric, results chan<- []QueryData, errs chan<- error) {
	for metric := range metrics {
		log.WithFields(log.Fields{
			"metric": metric.Name,
			"worker": workerID,
		}).Debug("Handling Topic")

		res, err := client.QueryMetric(cluster, metric, granularity, startTime, endTime)
		if err != nil {
			errs <- err
		} else {
			results <- res
		}
		wg.Done()

		log.WithFields(log.Fields{
			"metric": metric.Name,
			"worker": workerID,
		}).Debug("Handled Topic")
	}
	log.WithFields(log.Fields{
		"worker": workerID,
	}).Debug("Worker Done")
}

//QueryMetricWithAggs returns all the data points for a given metric, aggregated up to the given granularity, within the given window of time, and grouped by the given labels
func (client MetricsClient) QueryMetricWithAggs(cluster string, metric Metric, granularity Granularity, startTime time.Time, endTime time.Time, whitelistedLabels []MetricLabel) ([]QueryData, error) {
	query := Query{
		Filter:      NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
		Intervals:   []string{NewTimeInterval(startTime, endTime)},
		Aggreations: []Aggregation{NewMetricAgg(metric.Name)},
		Granularity: granularity.String(),
		GroupBy:     metric.GetValidLabels(whitelistedLabels),
		Limit:       DefaultQueryLimit,
	}

	response, err := client.SendQuery(queryPath, query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricAndType returns all the data points for a given metric and type, aggregated up to the given granularity, within the given window of time
func (client MetricsClient) QueryMetricAndType(cluster string, metric Metric, tye string, granularity Granularity, startTime time.Time, endTime time.Time) ([]QueryData, error) {
	query := Query{
		Filter:      NewFilterCollection(OpAnd, NewClusterFilter(cluster), NewTypeFilter(tye)),
		Intervals:   []string{NewTimeInterval(startTime, endTime)},
		Aggreations: []Aggregation{NewMetricAgg(metric.Name)},
		Granularity: granularity.String(),
		GroupBy:     metric.GetValidLabels([]MetricLabel{MetricLabelCluster, MetricLabelType}),
		Limit:       DefaultQueryLimit,
	}

	response, err := client.SendQuery(queryPath, query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricAndTopic returns all the data points for a given metric and topic, aggregated up to the given granularity, within the given window of time
func (client MetricsClient) QueryMetricAndTopic(cluster string, metric Metric, topic string, granularity Granularity, startTime time.Time, endTime time.Time, includePartitions bool) ([]QueryData, error) {
	if topic == "*" || strings.ToLower(topic) == "all" {
		return client.QueryMetricForAllTopics(cluster, metric, granularity, startTime, endTime, includePartitions, nil)
	}

	labels := []MetricLabel{MetricLabelCluster, MetricLabelTopic}
	if includePartitions {
		labels = append(labels, MetricLabelPartition)
	}
	query := Query{
		Filter:      NewFilterCollection(OpAnd, NewClusterFilter(cluster), NewTopicFilter(topic)),
		Intervals:   []string{NewTimeInterval(startTime, endTime)},
		Aggreations: []Aggregation{NewMetricAgg(metric.Name)},
		Granularity: granularity.String(),
		GroupBy:     metric.GetValidLabels(labels),
		Limit:       DefaultQueryLimit,
	}

	response, err := client.SendQuery(queryPath, query)
	for i, r := range response.Data {
		d := r
		d.Metric = metric.Name
		response.Data[i] = d
	}
	return response.Data, err
}

//QueryMetricAndTopicWorker returns all the data points, fetched in parallel, for a given metric and topics, aggregated up to the given granularity, within the given window of time
func (client MetricsClient) queryMetricAndTopicWorker(cluster string, metric Metric, granularity Granularity, startTime time.Time, endTime time.Time, includePartitions bool, workerID int, wg *sync.WaitGroup, topics <-chan string, results chan<- []QueryData, errs chan<- error) {
	for topic := range topics {
		log.WithFields(log.Fields{
			"topic":  topic,
			"worker": workerID,
		}).Debug("Handling Topic")

		res, err := client.QueryMetricAndTopic(cluster, metric, topic, granularity, startTime, endTime, includePartitions)
		if err != nil {
			errs <- err
		} else {
			results <- res
		}
		wg.Done()

		log.WithFields(log.Fields{
			"topic":  topic,
			"worker": workerID,
		}).Debug("Handled Topic")
	}
	log.WithFields(log.Fields{
		"worker": workerID,
	}).Debug("Worker Done")
}

//QueryMetricAndTopics returns all the data points, fetched in parallel, for a given metric and topics, aggregated up to the given granularity, within the given window of time
func (client MetricsClient) QueryMetricAndTopics(cluster string, metric Metric, topics []string, granularity Granularity, startTime time.Time, endTime time.Time, includePartitions bool, blacklistedTopics []string) ([]QueryData, error) {
	hasAll := false
	for _, t := range topics {
		if t == "*" || strings.ToLower(t) == "all" {
			hasAll = true
			break
		}
	}
	if hasAll {
		return client.QueryMetricForAllTopics(cluster, metric, granularity, startTime, endTime, includePartitions, blacklistedTopics)
	}
	topicChan := make(chan string, len(topics))
	resultsChan := make(chan []QueryData, len(topics))
	errorsChan := make(chan error, len(topics))
	waitGroup := sync.WaitGroup{}

	log.Debug("Starting up routines")
	for id := 0; id < client.apiContext.MaxWorkers; id++ {
		go client.queryMetricAndTopicWorker(cluster, metric, granularity, startTime, endTime, includePartitions, id, &waitGroup, topicChan, resultsChan, errorsChan)
	}

	log.Debug("Sending Topics")
	for _, topic := range topics {
		log.Debug("Sending Topic: " + topic)
		topicChan <- topic
		waitGroup.Add(1)
	}
	log.Debug("Done Sending Topics. Closing Channel")
	close(topicChan)

	waitGroup.Wait()
	close(resultsChan)
	close(errorsChan)

	log.Debug("Processing Errors")
	var err error
	finalErrors := []string{}
	for err := range errorsChan {
		finalErrors = append(finalErrors, err.Error())
	}
	if len(finalErrors) > 0 {
		err = errors.New(strings.Join(finalErrors, "\n\n"))
	} else {
		err = nil
	}

	log.Debug("Processing Results")
	queryData := []QueryData{}
	for res := range resultsChan {
		queryData = append(queryData, res...)
	}

	return queryData, err
}

//QueryMetricForAllTopics returns all the data points, fetched in parallel, for a given metric and all available topics (As returned by GetTopicsForMetric), aggregated up to the given granularity, within the given window of time
func (client MetricsClient) QueryMetricForAllTopics(cluster string, metric Metric, granularity Granularity, startTime time.Time, endTime time.Time, includePartitions bool, blacklistedTopics []string) ([]QueryData, error) {
	topics, err := client.GetTopicsForMetric(cluster, metric, startTime, endTime)
	if err != nil {
		return nil, err
	}

	finalTopics := []string{}

OUTER:
	for _, t := range topics {
		if blacklistedTopics != nil {
			for _, bt := range blacklistedTopics {
				if t == bt {
					continue OUTER
				}
			}
		}
		finalTopics = append(finalTopics, t)
	}

	log.WithFields(log.Fields{
		"topics": finalTopics,
		"metric": metric,
	}).Debug("Getting Results for All topics")

	return client.QueryMetricAndTopics(cluster, metric, finalTopics, granularity, startTime, endTime, includePartitions, blacklistedTopics)
}

//SendGet send a HTTP GET request to the metrics API at the given path
func (client MetricsClient) SendGet(path string) ([]byte, error) {
	if log.IsLevelEnabled(log.InfoLevel) {
		log.WithFields(log.Fields{
			"path": path,
		}).Debug("Sending GET Request")
	}
	res, err := client.sendReq("GET", path, nil)
	if log.IsLevelEnabled(log.TraceLevel) {
		log.WithFields(log.Fields{
			"path":   path,
			"result": string(res),
		}).Trace("Received GET Response")
	}
	return res, err
}

//SendPost sends a HTTP POST request to the metrics API at the given path with the given Query as the post body
func (client MetricsClient) SendPost(path string, query Query) ([]byte, error) {
	err := query.Validate()
	if err != nil {
		return nil, err
	}

	jsonQuery, err := query.ToJSON()
	if err != nil {
		return nil, err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		log.WithFields(log.Fields{
			"path":  path,
			"query": string(jsonQuery),
		}).Trace("Sending POST Request")
	}
	res, err := client.sendReq("POST", path, jsonQuery)
	if log.IsLevelEnabled(log.TraceLevel) {
		log.WithFields(log.Fields{
			"path":   path,
			"query":  string(jsonQuery),
			"result": string(res),
		}).Trace("Recieved POST Response")
	}
	return res, err
}

//SendQuery sends a HTTP POST request to the metrics API at the given path with the given Query as the post body and Unmarshals the resulted JSON into a QueryResponse
func (client MetricsClient) SendQuery(path string, query Query) (QueryResponse, error) {
	result, err := client.SendPost(path, query)

	if err != nil {
		return QueryResponse{}, err
	}

	response := QueryResponse{}
	json.Unmarshal(result, &response)

	return response, nil
}

func (client MetricsClient) sendReq(method string, path string, body []byte) ([]byte, error) {
	log.WithFields(log.Fields{
		"method": method,
		"path":   path,
	}).Trace("Sending API Request")

	endpoint := client.apiContext.BaseURL + path
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(client.apiContext.APIKey, client.apiContext.APISecret)

	req.Header.Add("Content-Type", "application/json")
	if client.httpContext.UserAgent != "" {
		req.Header.Add("User-Agent", client.httpContext.UserAgent)
	}

	for header, value := range client.httpContext.HTTPHeaders {
		req.Header.Add(header, value)
	}

	res, err := client.httpClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error returned from HTTP Request")
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err.Error(),
			"statusCode":    res.StatusCode,
			"statusMessage": res.Status,
		}).Error("Error returned from HTTP Request")
		return nil, err
	}

	if res.StatusCode != 200 {
		response := ErrorResponse{}
		json.Unmarshal(resBody, &response)

		errorMsg := fmt.Sprintf("Received status code %d instead of 200 for %s on %s", res.StatusCode, method, endpoint)
		log.WithFields(log.Fields{
			"statusCode":    res.StatusCode,
			"statusMessage": res.Status,
			"reqBody":       string(body),
			"errorResponse": response,
		}).Error(errorMsg)
		return nil, errors.New(errorMsg)
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		log.WithFields(log.Fields{
			"method":  method,
			"path":    path,
			"results": string(resBody),
		}).Trace("Api Request Results")
	}

	return resBody, nil
}
