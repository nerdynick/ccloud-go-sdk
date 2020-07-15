package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	availableMetrics = &cobra.Command{
		Use:   "metrics",
		Short: "List currently available metrics",
		RunE:  runE(&AvailableMetrics{}),
	}
)

type AvailableMetrics struct {
	Results []ccloudmetrics.AvailableMetric
}

func (am *AvailableMetrics) req(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient) error {
	res, err := client.GetAvailableMetrics()
	am.Results = res
	log.WithFields(log.Fields{
		"result": res,
		"err":    err,
	}).Info("Fetched Available Metrics")

	return err
}
func (am AvailableMetrics) outputPlain() error {
	log.WithFields(log.Fields{
		"result": am.Results,
	}).Info("Printing Plain Output")

	for _, metric := range am.Results {
		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, label.Name)
		}

		fmt.Printf("Name:      %s\n", metric.Name)
		fmt.Printf("Desc:      %s\n", metric.Desc)
		fmt.Printf("Type:      %s\n", metric.Type)
		fmt.Printf("LifeCycle: %s\n", metric.LifecycleStage)
		fmt.Printf("Labels:    %s\n", strings.Join(labels, ","))
		fmt.Println()
	}
	return nil
}
func (am AvailableMetrics) outputJSON(encoder *json.Encoder) error {
	return encoder.Encode(am.Results)
}
func (am AvailableMetrics) outputCSV(writer *csv.Writer) error {
	for _, metric := range am.Results {
		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, label.Name)
		}
		err := writer.Write([]string{
			metric.Name,
			metric.Desc,
			metric.Type,
			metric.LifecycleStage,
			strings.Join(labels, ";"),
		})
		if err != nil {
			return nil
		}
	}
	return nil
}

func init() {
	listCmd.AddCommand(availableMetrics)
}
