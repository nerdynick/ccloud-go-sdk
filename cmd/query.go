package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query metrics back from the API",
}

type Query struct {
	Results []ccloudmetrics.QueryData
	request func(*cobra.Command, []string, ccloudmetrics.MetricsClient, time.Time, time.Time) ([]ccloudmetrics.QueryData, error)
}

func (q *Query) req(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient) error {
	res, err := q.request(cmd, args, client, context.getStartTime(), context.getEndTime())
	q.Results = res
	outputErrs := []error{}
	if err != nil {
		outputErrs = append(outputErrs, err)
	}

	if res != nil {
		if outputCSV {
			writer := csv.NewWriter(os.Stdout)
			defer writer.Flush()
			err := q.outputCSV(writer)
			if err != nil {
				outputErrs = append(outputErrs, err)
			}

		}
		if outputJSON {
			encoder := json.NewEncoder(os.Stdout)
			err := q.outputJSON(encoder)
			if err != nil {
				outputErrs = append(outputErrs, err)
			}
		}
		if !outputCSV && !outputJSON {
			err := q.outputPlain()
			if err != nil {
				outputErrs = append(outputErrs, err)
			}
		}
	}

	finalErrors := []string{}
	for _, err := range outputErrs {
		finalErrors = append(finalErrors, err.Error())
	}
	if len(finalErrors) > 0 {
		return errors.New(strings.Join(finalErrors, "\n\n"))
	}
	return nil
}

func (q *Query) outputPlain() error {
	log.WithFields(log.Fields{
		"result": q.Results,
	}).Info("Printing Plain Output")

	for _, result := range q.Results {
		fmt.Printf("Timestamp: %s\n", result.Timestamp)
		fmt.Printf("Type:      %s\n", result.Type)
		fmt.Printf("Cluster:   %s\n", result.Cluster)
		fmt.Printf("Topic:     %s\n", result.Topic)
		fmt.Printf("Partition: %s\n", result.Partition)
		fmt.Printf("Value:     %f\n", result.Value)
		fmt.Println()
	}
	return nil
}
func (q Query) outputJSON(encoder *json.Encoder) error {
	return encoder.Encode(q.Results)
}
func (q Query) outputCSV(writer *csv.Writer) error {
	for _, result := range q.Results {
		err := writer.Write([]string{
			result.Timestamp,
			result.Type,
			result.Cluster,
			result.Topic,
			result.Partition,
			fmt.Sprintf("%f", result.Value),
		})
		if err != nil {
			return nil
		}
	}
	return nil
}

func init() {
	queryCmd.PersistentFlags().StringVarP(&context.Cluster, "cluster", "c", "", "Confluent Cloud Cluster ID")
	queryCmd.MarkPersistentFlagRequired("cluster")

	queryCmd.PersistentFlags().StringVar(&context.StartTime, "start", time.Now().Add(time.Duration(-1)*time.Hour).Format(ccloudmetrics.TimeFormatStr), "Start Time in the format of "+ccloudmetrics.TimeFormatStr)
	queryCmd.PersistentFlags().StringVar(&context.EndTime, "end", time.Now().Format(ccloudmetrics.TimeFormatStr), "End Time in the format of "+ccloudmetrics.TimeFormatStr)
	queryCmd.PersistentFlags().IntVar(&context.LastXmin, "last", 0, "Instead of using start/end time. Query for the last X mins")

	queryCmd.PersistentFlags().StringVar(&context.Granularity, "gran", string(ccloudmetrics.GranularityOneHour), "Granularity of Metrics. Options are: "+strings.Join(ccloudmetrics.AvailableGranularities, ", "))

	rootCmd.AddCommand(queryCmd)
}
