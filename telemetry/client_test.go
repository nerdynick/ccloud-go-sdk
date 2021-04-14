package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiPaths(t *testing.T) {
	assert := assert.New(t)

	apiClient := New("", "")

	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/query", apiPathsQuery.format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/attributes", apiPathsAttributes.format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/descriptors", apiPathsDescriptor.format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/descriptors/metrics", apiPathsDescriptorMetrics.format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/descriptors/resources", apiPathsDescriptorResources.format(apiClient, 1))
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	apiClient := New("apikey", "apisec")

	assert.Equal("apikey", apiClient.Context.APIKey)
	assert.Equal("apisec", apiClient.Context.APISecret)
	assert.Equal(DefaultBaseURL, apiClient.BaseURL)
	assert.Equal(DefaultQueryLimit, apiClient.PageLimit)
	assert.Equal(DatasetCloud, apiClient.DataSet)
	assert.Equal(DefaultMaxWorkers, apiClient.MaxWorkers)
}
