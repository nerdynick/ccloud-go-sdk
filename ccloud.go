package ccloudmetrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	descriptorPath string = "/v1/metrics/cloud/descriptors"
	attributesPath string = "/v1/metrics/cloud/attributes"
	availablePath  string = "/v1/metrics/cloud/available"
	queryPath      string = "/v1/metrics/cloud/query"

	opAnd string = "AND"
	opEq  string = "EQ"
	opOr  string = "OR"

	AggSum string = "SUM"

	MetricLabelCluster    string = "metric.label.cluster_id"
	MetricLabelTopic      string = "metric.label.topic"
	MetricLabelType       string = "metric.label.type"
	GranularityOneMin     string = "PT1M"
	GranularityFiveMin    string = "PT5M"
	GranularityFifteenMin string = "PT15M"
	GranularityThirtyMin  string = "PT30M"
	GranularityOneHour    string = "PT1H"
	GranularityAll        string = "ALL"
	LifecycleStagePreview string = "PREVIEW"
	LifecycleStageGeneral string = "GENERAL_AVAILABILITY"
)

type ApiContext struct {
	APIKey    string
	APISecret string
	BaseURL   string
}

type HTTPContext struct {
	RequestTimeout int
	UserAgent      string
	HTTPHeaders    map[string]string
}

type MetricsClient struct {
	apiContext  ApiContext
	httpContext HTTPContext
	httpClient  http.Client
}

func (client MetricsClient) GetCurrentlyAvailableMetrics(cluster string) ([]CurrentlyAvailableMetric, error) {
	filter := Filter{
		Field: MetricLabelCluster,
		Op:    opEq,
		Value: cluster,
	}
	query := Query{
		Filter: FilterHeader{
			Op:      opAnd,
			Filters: []Filter{filter},
		},
	}
	result, err := client.SendQuery(availablePath, query)

	if err != nil {
		return nil, err
	}

	response := CurrentlyAvailableMetricResponse{}
	json.Unmarshal(result, &response)

	return *response.AvailableMetrics, nil
}

func (client MetricsClient) GetAvailableMetrics() ([]AvailableMetric, error) {
	result, err := client.SendGet(descriptorPath)

	if err != nil {
		return nil, err
	}

	response := AvailableMetricResponse{}
	json.Unmarshal(result, &response)

	return *response.AvailableMetrics, nil
}

func (client MetricsClient) ListTopicsForMetric(cluster string, metric string) ([]string, error) {
	query := Query{
		Filter: FilterHeader{
			Op: opAnd,
			Filters: []Filter{Filter{
				Field: MetricLabelCluster,
				Op:    opEq,
				Value: cluster,
			}},
		},
		GroupBy: []string{MetricLabelTopic},
	}
	result, err := client.SendQuery(availablePath, query)

	if err != nil {
		return nil, err
	}

	response := QueryResponse{}
	json.Unmarshal(result, &response)

	topics := make([]string, len(response.Data))
	for i := 0; i < len(response.Data); i++ {
		topics[i] = response.Data[i].Topic
	}

	return topics, nil
}

func (client MetricsClient) QueryMetric(cluster, metric string, topic string) ([]QueryData, error) {
	return nil, nil
}

func (client MetricsClient) SendGet(path string) ([]byte, error) {
	return client.sendReq("GET", path, nil)
}

func (client MetricsClient) SendQuery(path string, query Query) ([]byte, error) {
	jsonQuery, err := query.ToJSON()
	if err != nil {
		panic(err)
	}
	return client.sendReq("POST", path, bytes.NewBuffer(jsonQuery))
}

func (client MetricsClient) sendReq(method string, path string, body io.Reader) ([]byte, error) {
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
		fmt.Printf(err.Error())
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := fmt.Sprintf("Received status code %d instead of 200 for %s on %s", res.StatusCode, method, endpoint)
		return nil, errors.New(errorMsg)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return resBody, nil
}

func NewClient(context *ApiContext, httpContext *HTTPContext) MetricsClient {
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(httpContext.RequestTimeout),
	}
	client := MetricsClient{
		httpClient:  httpClient,
		apiContext:  *context,
		httpContext: *httpContext,
	}

	return client
}
