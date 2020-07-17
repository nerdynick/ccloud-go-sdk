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
			return client.QueryMetricAndTopic(context.Cluster, context.getMetric(), context.Topic, context.getGranularity(), sTime, eTime, context.IncludePartitions)
		},
	}),
}

var topicsQueryCmd = &cobra.Command{
	Use:   "topics",
	Short: "Query a collection of topics for a particular metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetricAndTopics(context.Cluster, context.getMetric(), context.Topics, context.getGranularity(), sTime, eTime, context.IncludePartitions)
		},
	}),
}

var topicsAllQueryCmd = &cobra.Command{
	Use:   "all",
	Short: "Query all topics for a particular metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, sTime time.Time, eTime time.Time) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetricForAllTopics(context.Cluster, context.getMetric(), context.getGranularity(), sTime, eTime, context.IncludePartitions, context.BlacklistedTopics)
		},
	}),
}

func init() {
	topicQueryCmd.Flags().StringVarP(&context.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicQueryCmd.MarkFlagRequired("metric")
	topicQueryCmd.Flags().StringVar(&context.Topic, "topic", "", "Topic to query metric for")
	topicQueryCmd.MarkFlagRequired("topic")
	topicQueryCmd.Flags().BoolVar(&context.IncludePartitions, "partitions", false, "Should results be aggrigated to the parition or just to the topic")
	queryCmd.AddCommand(topicQueryCmd)

	topicsAllQueryCmd.Flags().StringVarP(&context.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicsAllQueryCmd.MarkFlagRequired("metric")
	topicsAllQueryCmd.Flags().StringArrayVar(&context.BlacklistedTopics, "blacklist", []string{}, "List of Topics to blacklist from getting fetch")
	topicsQueryCmd.AddCommand(topicsAllQueryCmd)

	topicsQueryCmd.Flags().StringVarP(&context.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicsQueryCmd.MarkFlagRequired("metric")
	topicsQueryCmd.Flags().StringArrayVar(&context.Topics, "topics", []string{}, "List of Topics to query for")
	topicsQueryCmd.MarkFlagRequired("topics")
	queryCmd.AddCommand(topicsQueryCmd)
}
