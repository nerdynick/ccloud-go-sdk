package cmd

import (
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var ()

var metricQueryCmd = &cobra.Command{
	Use:   "metric",
	Short: "Query results for a given metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetric(context.Cluster, context.getMetric(), context.getGranularity(), sTime, eTime)
		},
	}),
}

var metricsQueryCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Query results for a given metrics",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetrics(context.Cluster, context.getMetrics(), context.getGranularity(), sTime, eTime)
		},
	}),
}

func init() {
	metricQueryCmd.Flags().StringVarP(&context.Metric, "metric", "m", "", "Metric to fetch available topics for")
	metricQueryCmd.MarkFlagRequired("metric")
	queryCmd.AddCommand(metricQueryCmd)

	metricsQueryCmd.Flags().StringArrayVarP(&context.Metrics, "metric", "m", []string{}, "Metric to fetch available topics for. Can be used multipule times to provide more multipule metrics")
	metricsQueryCmd.MarkFlagRequired("metric")
	queryCmd.AddCommand(metricsQueryCmd)

}
