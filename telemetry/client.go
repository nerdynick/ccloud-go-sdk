package telemetry

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/nerdynick/ccloud-go-sdk/client"
)

const (
	//DefaultQueryLimit the default query limit for results
	DefaultQueryLimit int = 1000
	//DefaultBaseURL is the default Domain and Protocol for quering against the Metrics API
	DefaultBaseURL string = "https://api.telemetry.confluent.cloud"
	//DefaultMaxWorkers controls the max number of workers in a given Worker Pool that will be spawned
	DefaultMaxWorkers int = 5

	apiPathsQuery               apiPaths = "%s/v%d/metrics/%s/query"
	apiPathsAttributes          apiPaths = "%s/v%d/metrics/%s/attributes"
	apiPathsDescriptor          apiPaths = "%s/v%d/metrics/%s/descriptors"
	apiPathsDescriptorMetrics   apiPaths = "%s/v%d/metrics/%s/descriptors/metrics"
	apiPathsDescriptorResources apiPaths = "%s/v%d/metrics/%s/descriptors/resources"

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

func (p apiPaths) format(telemetryClient TelemetryClient, version int8) string {
	return fmt.Sprintf(string(p), telemetryClient.BaseURL, version, string(telemetryClient.DataSet))
}

//TelemetryClient is the SDK Client for making REST calls to the Confluent Metrics API
type TelemetryClient struct {
	client.Client
	PageLimit  int
	DataSet    Dataset
	BaseURL    string
	MaxWorkers int
}

//New Used to create a new MetricsClient from the given minimal set of properties
func New(apiKey string, apiSecret string) TelemetryClient {
	return TelemetryClient{
		Client:     client.New(apiKey, apiSecret),
		DataSet:    DatasetCloud,
		PageLimit:  DefaultQueryLimit,
		BaseURL:    DefaultBaseURL,
		MaxWorkers: DefaultMaxWorkers,
	}
}
