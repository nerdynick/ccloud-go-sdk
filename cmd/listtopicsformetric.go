package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var topicsForMetric = &cobra.Command{
	Use:   "topics",
	Short: "List all available topics for a given metric",
	RunE:  runE(&TopicsForMetric{}),
}

type TopicsForMetric struct {
	Results []string
}

func (am *TopicsForMetric) req(cmd *cobra.Command, args []string, context RequestContext, client ccloudmetrics.MetricsClient) (bool, error) {
	res, err := client.GetTopicsForMetric(context.Cluster, context.getMetric(), context.getStartTime(), context.getEndTime())

	am.Results = res
	log.WithFields(log.Fields{
		"result":  res,
		"err":     err,
		"context": context,
	}).Info("Fetched Available Topics for Metric")

	return (len(res) > 0), err
}
func (am TopicsForMetric) outputPlain() error {
	log.WithFields(log.Fields{
		"result": am.Results,
	}).Info("Printing Plain Output")

	for _, topic := range am.Results {
		fmt.Printf("Topic: %s\n", topic)
		fmt.Println()
	}
	return nil
}
func (am TopicsForMetric) outputJSON(encoder *json.Encoder) error {
	return encoder.Encode(am.Results)
}
func (am TopicsForMetric) outputCSV(writer *csv.Writer) error {
	for _, topic := range am.Results {
		err := writer.Write([]string{topic})
		if err != nil {
			return nil
		}
	}
	return nil
}

func init() {
	topicsForMetric.Flags().StringVarP(&requestcontext.Cluster, "cluster", "c", "", "Confluent Cloud Cluster ID")
	topicsForMetric.MarkFlagRequired("cluster")

	topicsForMetric.Flags().StringVarP(&requestcontext.Metric, "metric", "m", "", "Metric to fetch available topics for")
	topicsForMetric.MarkFlagRequired("metric")

	topicsForMetric.Flags().StringVar(&requestcontext.StartTime, "start", time.Now().Add(time.Duration(-1)*time.Hour).Format(ccloudmetrics.TimeFormatStr), "Start Time in the format of "+ccloudmetrics.TimeFormatStr)
	topicsForMetric.Flags().StringVar(&requestcontext.EndTime, "end", time.Now().Format(ccloudmetrics.TimeFormatStr), "End Time in the format of "+ccloudmetrics.TimeFormatStr)
	listCmd.AddCommand(topicsForMetric)
}
