package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiPaths(t *testing.T) {
	assert := assert.New(t)

	apiClient := New("", "")

	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/query", APIPathQuery.Format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/attributes", APIPathAttributes.Format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/descriptors", APIPathDescriptor.Format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/descriptors/metrics", APIPathDescriptorMetrics.Format(apiClient, 1))
	assert.Equal(DefaultBaseURL+"/v1/metrics/cloud/descriptors/resources", APIPathDescriptorResources.Format(apiClient, 1))
}
