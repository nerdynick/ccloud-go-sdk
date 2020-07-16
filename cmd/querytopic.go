package cmd

import (
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var topicQueryCmd = &cobra.Command{
	Use:   "topic",
	Short: "Query a topic for a particular metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetricAndTopic(cluster, metric, topic, granularity, sTime, eTime, includePartitions)
		},
	}),
}

var (
	topic             string
	includePartitions bool
)

func init() {
	topicQueryCmd.Flags().StringVar(&topic, "topic", "", "Topic to query metric for")
	topicQueryCmd.MarkFlagRequired("topic")
	topicQueryCmd.Flags().BoolVar(&includePartitions, "partitions", false, "Should results be aggrigated to the parition or just to the topic")
	queryCmd.AddCommand(topicQueryCmd)
}
