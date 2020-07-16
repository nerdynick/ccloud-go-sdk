package ccloudmetrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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

//APIContext is the Contextual set of configs for all Metrics API calls
type APIContext struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	MaxWorkers int
}

//HTTPContext is the Contextual set of configs for the HTTP Client making the calls to the Metrics API
type HTTPContext struct {
	RequestTimeout      int
	UserAgent           string
	HTTPHeaders         map[string]string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

//MetricsClient is the SDK Client for making REST calls to the Confluent Metrics API
type MetricsClient struct {
	apiContext  APIContext
	httpContext HTTPContext
	httpClient  http.Client
}

//NewClientFromContext Used to create a new MetricsClient from the given contexts
func NewClientFromContext(context *APIContext, httpContext *HTTPContext) MetricsClient {
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(httpContext.RequestTimeout),
		Transport: &http.Transport{
			MaxIdleConns:        httpContext.MaxIdleConns,
			MaxIdleConnsPerHost: httpContext.MaxIdleConnsPerHost,
		},
	}
	client := MetricsClient{
		httpClient:  httpClient,
		apiContext:  *context,
		httpContext: *httpContext,
	}

	return client
}

//NewClientMinimal Used to create a new MetricsClient from the given minimal set of properties
func NewClientMinimal(apiKey string, apiSecret string) MetricsClient {
	return NewClientWithDefaults(apiKey, apiSecret, nil)
}

//NewClientWithDefaults Used to create a new MetricsClient from the minimal set of properties and using defaults where appropriate
func NewClientWithDefaults(apiKey string, apiSecret string, extraHeaders map[string]string) MetricsClient {
	return NewClient(apiKey, apiSecret, DefaultBaseURL, DefaultRequestTimeout, DefaultUserAgent, DefaultMaxWorkers, extraHeaders)
}

//NewClient Used to create a new MetricsClient from the full set of properties
func NewClient(apiKey string, apiSecret string, baseURL string, requestTime int, userAgent string, maxWorkers int, extraHeaders map[string]string) MetricsClient {
	return NewClientFromContext(&APIContext{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    strOrDefault(baseURL, DefaultBaseURL),
		MaxWorkers: intOrDefault(maxWorkers, DefaultMaxWorkers),
	}, &HTTPContext{
		RequestTimeout: intOrDefault(requestTime, DefaultRequestTimeout),
		UserAgent:      strOrDefault(userAgent, DefaultUserAgent),
		HTTPHeaders:    extraHeaders,
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

func (client MetricsClient) GetAvailableMetrics() ([]AvailableMetric, error) {
	result, err := client.SendGet(descriptorPath)

	if err != nil {
		return nil, err
	}

	response := AvailableMetricResponse{}
	json.Unmarshal(result, &response)

	return response.AvailableMetrics, nil
}

func (client MetricsClient) GetCurrentlyAvailableMetrics(cluster string) ([]CurrentlyAvailableMetric, error) {
	query := Query{
		Filter: NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
	}
	result, err := client.SendPost(availablePath, query)

	if err != nil {
		return nil, err
	}

	response := CurrentlyAvailableMetricResponse{}
	json.Unmarshal(result, &response)

	return response.AvailableMetrics, nil
}

func (client MetricsClient) GetTopicsForMetric(cluster string, metric string, startTime time.Time, endTime time.Time) ([]string, error) {
	query := Query{
		Filter:    NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
		GroupBy:   []string{MetricLabelTopic},
		Intervals: []string{NewTimeInterval(startTime, endTime)},
		Metric:    metric,
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

func (client MetricsClient) QueryMetric(cluster string, metric string, granularity string, startTime time.Time, endTime time.Time) ([]QueryData, error) {
	query := Query{
		Filter:      NewFilterCollection(OpAnd, NewClusterFilter(cluster)),
		Intervals:   []string{NewTimeInterval(startTime, endTime)},
		Aggreations: []Aggregation{NewMetricAgg(metric)},
		Granularity: granularity,
		GroupBy:     []string{MetricLabelCluster},
		Limit:       DefaultQueryLimit,
	}

	response, err := client.SendQuery(attributesPath, query)
	return response.Data, err
}

func (client MetricsClient) QueryMetricAndTopic(cluster string, metric string, topic string, granularity string, startTime time.Time, endTime time.Time, includePartitions bool) ([]QueryData, error) {
	query := Query{
		Filter:      NewFilterCollection(OpAnd, NewClusterFilter(cluster), NewTopicFilter(topic)),
		Intervals:   []string{NewTimeInterval(startTime, endTime)},
		Aggreations: []Aggregation{NewMetricAgg(metric)},
		Granularity: granularity,
		GroupBy:     []string{MetricLabelCluster, MetricLabelTopic},
		Limit:       DefaultQueryLimit,
	}

	if includePartitions {
		query.GroupBy = append(query.GroupBy, MetricLabelPartition)
	}

	response, err := client.SendQuery(attributesPath, query)
	return response.Data, err
}

func (client MetricsClient) QueryMetricAndTopicWorker(cluster string, metric string, granularity string, startTime time.Time, endTime time.Time, includePartitions bool, workerId int, topics <-chan string, results chan<- []QueryData, errs chan<- error) {
	for topic := range topics {
		res, err := client.QueryMetricAndTopic(cluster, metric, topic, granularity, startTime, endTime, includePartitions)
		if err != nil {
			errs <- err
		} else {
			results <- res
		}

	}
}

func (client MetricsClient) QueryMetricForTopics(cluster, metric string, topics []string, granularity string, startTime time.Time, endTime time.Time, includePartitions bool) ([]QueryData, error) {
	topicChan := make(chan string, client.apiContext.MaxWorkers)
	resultsChan := make(chan []QueryData, client.apiContext.MaxWorkers)
	errorsChan := make(chan error, client.apiContext.MaxWorkers)

	for id := 0; id < client.apiContext.MaxWorkers; id++ {
		go client.QueryMetricAndTopicWorker(cluster, metric, granularity, startTime, endTime, includePartitions, id, topicChan, resultsChan, errorsChan)
	}

	for _, topic := range topics {
		topicChan <- topic
	}
	close(topicChan)

	finalErrors := []string{}
	for err := range errorsChan {
		finalErrors = append(finalErrors, err.Error())
	}
	if len(finalErrors) > 0 {
		return nil, errors.New(strings.Join(finalErrors, "\n\n"))
	}

	queryData := []QueryData{}
	for res := range resultsChan {
		queryData = append(queryData, res...)
	}

	return queryData, nil
}

func (client MetricsClient) QueryMetricForAllTopics(cluster, metric string, topics []string, granularity string, startTime time.Time, endTime time.Time, includePartitions bool, blacklistedTopics []string) ([]QueryData, error) {
	topics, err := client.GetTopicsForMetric(cluster, metric, startTime, endTime)
	if err != nil {
		return nil, err
	}

	finalTopics := []string{}

OUTER:
	for _, t := range topics {
		for _, bt := range blacklistedTopics {
			if t == bt {
				continue OUTER
			}
		}
		finalTopics = append(finalTopics, t)
	}

	return client.QueryMetricForTopics(cluster, metric, finalTopics, granularity, startTime, endTime, includePartitions)
}

func (client MetricsClient) SendGet(path string) ([]byte, error) {
	if log.IsLevelEnabled(log.InfoLevel) {
		log.WithFields(log.Fields{
			"path": path,
		}).Info("Sending GET Request")
	}
	return client.sendReq("GET", path, nil)
}

func (client MetricsClient) SendPost(path string, query Query) ([]byte, error) {
	jsonQuery, err := query.ToJSON()
	if err != nil {
		panic(err)
	}
	if log.IsLevelEnabled(log.InfoLevel) {
		log.WithFields(log.Fields{
			"path":  path,
			"query": string(jsonQuery),
		}).Info("Sending POST Request")
	}
	return client.sendReq("POST", path, bytes.NewBuffer(jsonQuery))
}

func (client MetricsClient) SendQuery(path string, query Query) (QueryResponse, error) {
	result, err := client.SendPost(path, query)

	if err != nil {
		return QueryResponse{}, err
	}

	response := QueryResponse{}
	json.Unmarshal(result, &response)

	return response, nil
}

func (client MetricsClient) sendReq(method string, path string, body io.Reader) ([]byte, error) {
	log.WithFields(log.Fields{
		"method": method,
		"path":   path,
	}).Trace("Sending API Request")

	endpoint := client.apiContext.BaseURL + path
	req, err := http.NewRequest(method, endpoint, body)
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
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error returned from HTTP Request")
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		errorMsg := fmt.Sprintf("Received status code %d instead of 200 for %s on %s", res.StatusCode, method, endpoint)
		log.WithFields(log.Fields{
			"statusCode": res.StatusCode,
		}).Error(errorMsg)
		return nil, errors.New(errorMsg)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error returned from HTTP Request")
		return nil, err
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
