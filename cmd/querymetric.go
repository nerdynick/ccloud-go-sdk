package cmd

import (
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var (
	metrics []string
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

var metricsQueryCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Query results for a given metrics",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			results := []ccloudmetrics.QueryData{}
			for _, metric := range metrics {
				res, err := client.QueryMetric(cluster, metric, granularity, sTime, eTime)
				if res != nil {
					results = append(results, res...)
				}
				if err != nil {
					return results, err
				}
			}

			return results, nil
		},
	}),
}

func init() {
	metricQueryCmd.Flags().StringVarP(&metric, "metric", "m", "", "Metric to fetch available topics for")
	metricQueryCmd.MarkFlagRequired("metric")
	queryCmd.AddCommand(metricQueryCmd)

	metricsQueryCmd.Flags().StringArrayVarP(&metrics, "metric", "m", []string{}, "Metric to fetch available topics for. Can be used multipule times to provide more multipule metrics")
	metricsQueryCmd.MarkFlagRequired("metric")
	queryCmd.AddCommand(metricsQueryCmd)

}
