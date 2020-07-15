package cmd

import (
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List information from the API",
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	//Available Metrics Command
	listCmd.AddCommand(availableMetrics)

	//Topics for Metric Command
	addClusterFlag(topicsForMetric.Flags())
	addMetricFlag(topicsForMetric.Flags())
	topicsForMetric.Flags().StringVar(&startTime, "start", time.Now().Add(time.Duration(-1)*time.Hour).Format(ccloudmetrics.TimeFormatStr), "Start Time in the format of "+ccloudmetrics.TimeFormatStr)
	topicsForMetric.Flags().StringVar(&endTime, "end", time.Now().Format(ccloudmetrics.TimeFormatStr), "End Time in the format of "+ccloudmetrics.TimeFormatStr)
	listCmd.AddCommand(topicsForMetric)
}
