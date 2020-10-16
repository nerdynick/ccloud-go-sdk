package cmd

import (
	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var topicQueryCmd = &cobra.Command{
	Use:   "topic",
	Short: "Query a topic for a particular metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, context RequestContext) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetricAndTopic(context.Cluster, context.getMetric(), context.Topic, context.getGranularity(), context.getStartTime(), context.getEndTime(), context.IncludePartitions)
		},
	}),
}

var topicsQueryCmd = &cobra.Command{
	Use:   "topics",
	Short: "Query a collection of topics for a particular metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, context RequestContext) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetricAndTopics(context.Cluster, context.getMetric(), context.Topics, context.getGranularity(), context.getStartTime(), context.getEndTime(), context.IncludePartitions, context.BlacklistedTopics)
		},
	}),
}

var topicsAllQueryCmd = &cobra.Command{
	Use:   "all",
	Short: "Query all topics for a particular metric",
	RunE: runE(&Query{
		request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, context RequestContext) ([]ccloudmetrics.QueryData, error) {
			return client.QueryMetricForAllTopics(context.Cluster, context.getMetric(), context.getGranularity(), context.getStartTime(), context.getEndTime(), context.IncludePartitions, context.BlacklistedTopics)
		},
	}),
}

func init() {
	topicQueryCmd.Flags().StringVarP(&requestcontext.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicQueryCmd.MarkFlagRequired("metric")
	topicQueryCmd.Flags().StringVar(&requestcontext.Topic, "topic", "", "Topic to query metric for")
	topicQueryCmd.MarkFlagRequired("topic")
	topicQueryCmd.Flags().BoolVar(&requestcontext.IncludePartitions, "partitions", false, "Should results be aggrigated to the parition or just to the topic")
	queryCmd.AddCommand(topicQueryCmd)

	topicsAllQueryCmd.Flags().StringVarP(&requestcontext.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicsAllQueryCmd.MarkFlagRequired("metric")
	topicsAllQueryCmd.Flags().StringArrayVar(&requestcontext.BlacklistedTopics, "blacklist", []string{}, "List of Topics to blacklist from getting fetch")
	topicsQueryCmd.AddCommand(topicsAllQueryCmd)

	topicsQueryCmd.Flags().StringVarP(&requestcontext.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicsQueryCmd.MarkFlagRequired("metric")
	topicsQueryCmd.Flags().StringArrayVar(&requestcontext.Topics, "topics", []string{}, "List of Topics to query for")
	topicsQueryCmd.MarkFlagRequired("topics")
	queryCmd.AddCommand(topicsQueryCmd)
}
