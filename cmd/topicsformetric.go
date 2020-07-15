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

func (am *TopicsForMetric) req(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient) error {
	sTime, err := time.Parse(ccloudmetrics.TimeFormatStr, startTime)
	if err != nil {
		return nil
	}
	eTime, err := time.Parse(ccloudmetrics.TimeFormatStr, endTime)
	if err != nil {
		return nil
	}

	res, err := client.GetTopicsForMetric(cluster, metric, sTime, eTime)

	am.Results = res
	log.WithFields(log.Fields{
		"result":    res,
		"err":       err,
		"startTime": sTime,
		"endTime":   eTime,
		"metric":    metric,
		"cluster":   cluster,
	}).Info("Fetched Available Topics for Metric")

	return err
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
