package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	apiClient := New("apikey", "apisec")

	assert.Equal("apikey", apiClient.Context.APIKey)
	assert.Equal("apisec", apiClient.Context.APISecret)
	assert.Equal(DefaultBaseURL, apiClient.Context.BaseURL)
	assert.Equal(DefaultQueryLimit, apiClient.PageLimit)
	assert.Equal(DatasetCloud, apiClient.DataSet)
	assert.Equal(DefaultMaxWorkers, apiClient.MaxWorkers)
}
