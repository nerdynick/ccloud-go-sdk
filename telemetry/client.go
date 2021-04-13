package ccloudmetrics

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/client"
	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry/query"
)

const (
	//DefaultQueryLimit the default query limit for results
	DefaultQueryLimit int = 1000
	//DefaultBaseURL is the default Domain and Protocol for quering against the Metrics API
	DefaultBaseURL string = "https://api.telemetry.confluent.cloud"
	//DefaultMaxWorkers controls the max number of workers in a given Worker Pool that will be spawned
	DefaultMaxWorkers int = 5

	//DefaultAPIVersion API Version to call
	DefaultAPIVersion int8 = 1

	apiPathsQuery               apiPaths = "/v$d/metrics/$s/query"
	apiPathsAttributes          apiPaths = "/v$d/metrics/$s/attributes"
	apiPathsDescriptor          apiPaths = "/v$d/metrics/$s/descriptors"
	apiPathsDescriptorMetrics   apiPaths = "/v$d/metrics/$s/descriptors/metrics"
	apiPathsDescriptorResources apiPaths = "/v$d/metrics/$s/descriptors/resources"

	//DatasetCloud constant name for the CCloud dataset
	DatasetCloud  Dataset = "cloud"
	DatasetHosted Dataset = "hosted-monitoring"
)

var (
	//AvailableDatasets Constant for the currently known available Datasets
	AvailableDatasets []Dataset = []Dataset{
		DatasetCloud,
	}

	cJSON = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "cjson",
	}.Froze()
)

//Dataset struct to referece the selected dataset
type Dataset string

type apiPaths string

func (p apiPaths) format(version int8, dataset Dataset) string {
	return fmt.Sprintf(string(p), version, dataset)
}

//TelemetryClient is the SDK Client for making REST calls to the Confluent Metrics API
type TelemetryClient struct {
	client.Client
	PageLimit  int
	DataSet    Dataset
	APIKey     string
	APISecret  string
	BaseURL    string
	MaxWorkers int
}

//New Used to create a new MetricsClient from the given minimal set of properties
func New(apiKey string, apiSecret string) TelemetryClient {
	return TelemetryClient{
		Client: client.New(apiKey, apiSecret),
	}
}

func (client TelemetryClient) SendPostQuery(response interface{}, url string, q query.Query) error {
	err := q.Validate()
	if err != nil {
		return err
	}

	return client.SendPost(&response, url, q)
}
