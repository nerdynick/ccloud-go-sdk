package telemetry

import (
	"fmt"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/client"
)

const (
	APIPathQuery               TelemetryAPIPath = "metrics/%s/query"
	APIPathAttributes          TelemetryAPIPath = "metrics/%s/attributes"
	APIPathDescriptor          TelemetryAPIPath = "metrics/%s/descriptors"
	APIPathDescriptorMetrics   TelemetryAPIPath = "metrics/%s/descriptors/metrics"
	APIPathDescriptorResources TelemetryAPIPath = "metrics/%s/descriptors/resources"
)

type TelemetryAPIPath client.APIPath

func (p TelemetryAPIPath) Format(client TelemetryClient, apiVersion int8) string {
	return fmt.Sprintf(strings.Join([]string{
		client.Context.BaseURL,
		"v" + fmt.Sprint(apiVersion),
		string(p),
	}, "/"), client.DataSet)
}
