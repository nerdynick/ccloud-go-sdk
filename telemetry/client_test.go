package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiPaths(t *testing.T) {
	assert := assert.New(t)

	apiClient := New("", "")

	assert.Equal("/v1/metrics/cloud/query", apiPathsQuery.format(1, apiClient.DataSet))
	assert.Equal("/v1/metrics/cloud/attributes", apiPathsAttributes.format(1, apiClient.DataSet))
	assert.Equal("/v1/metrics/cloud/descriptors", apiPathsDescriptor.format(1, apiClient.DataSet))
	assert.Equal("/v1/metrics/cloud/descriptors/metrics", apiPathsDescriptorMetrics.format(1, apiClient.DataSet))
	assert.Equal("/v1/metrics/cloud/descriptors/resources", apiPathsDescriptorResources.format(1, apiClient.DataSet))
}
