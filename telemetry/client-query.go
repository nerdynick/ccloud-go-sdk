package telemetry

import (
	"github.com/nerdynick/ccloud-go-sdk/logging"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query"
	"go.uber.org/zap"
)

//PostQuery POST Query to the Telemetry API
func (client TelemetryClient) PostQuery(response interface{}, url string, q query.Query) error {
	if client.Log.Core().Enabled(logging.InfoLevel) {
		qJson, _ := q.ToJSON()
		client.Log.Info("Query - Posting",
			zap.String("URI", url),
			zap.ByteString("Query", qJson),
		)
	}

	err := q.Validate()
	if err != nil {
		return err
	}

	return client.Post(&response, url, q)
}
