package telemetry

import (
	"net/url"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/response"
)

func (client *TelemetryClient) SendDesc() (response.Metrics, error) {
	url := apiPathsDescriptor.format(*client, 1)
	response := response.Metrics{}

	err := client.Get(&response, url)
	return response, err
}

func (client *TelemetryClient) SendDescMetrics(resourceType labels.Resource) (response.Metrics, error) {
	url, _ := url.ParseRequestURI(apiPathsDescriptorMetrics.format(*client, 2))
	q := url.Query()
	q.Add("resource_type", resourceType.Key)
	url.RawQuery = q.Encode()

	response := response.Metrics{}

	err := client.Get(&response, url.String())
	return response, err
}

func (client *TelemetryClient) SendDescResources() (response.Resources, error) {
	response := response.Resources{}
	url := apiPathsDescriptorResources.format(*client, 2)
	err := client.Get(&response, url)

	return response, err
}

//GetAvailableMetrics returns a collection of all the available metrics and their supported labels among other important meta data for Kafka Clusters
func (client *TelemetryClient) GetAvailableMetrics() ([]metric.Metric, error) {
	response, err := client.SendDesc()
	if err != nil {
		return nil, err
	}
	return response.AvailableMetrics, err
}

//GetAvailableMetricsForResource returns a collection of all the available metrics and their supported labels among other important meta data for a given resource type
// This is also a Preview V2 API feature and may be subject to breakage and/or change at any moment
func (client *TelemetryClient) GetAvailableMetricsForResource(resourceType labels.Resource) ([]metric.Metric, error) {
	response, err := client.SendDescMetrics(resourceType)
	if err != nil {
		return nil, err
	}
	return response.AvailableMetrics, err
}

//GetAvailableResources returns a collection of all the available metrics and their supported labels among other important meta data.
// This is also a Preview V2 API feature and may be subject to breakage and/or change at any moment
func (client *TelemetryClient) GetAvailableResources() ([]response.ResourceType, error) {
	response, err := client.SendDescResources()

	if err != nil {
		return nil, err
	}
	return response.ResourceTypes, err
}
