package cmd

import (
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var metricQueryCmd = &cobra.Command{
	Use:   "metric",
	Short: "Query results for a given metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetric(cluster, metric, granularity, sTime, eTime)
		},
	}),
}

func init() {
	queryCmd.AddCommand(metricQueryCmd)
}
