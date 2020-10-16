package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	queryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query metrics back from the API",
	}
)

type Query struct {
	Results []ccloudmetrics.QueryData
	request func(*cobra.Command, []string, ccloudmetrics.MetricsClient, RequestContext) ([]ccloudmetrics.QueryData, error)
}

func (q *Query) req(cmd *cobra.Command, args []string, context RequestContext, client ccloudmetrics.MetricsClient) (bool, error) {
	res, err := q.request(cmd, args, client, context)
	q.Results = res
	log.WithField("results", res).Info("Query Results")
	return (len(res) > 0), err
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
	cobra.OnInitialize(queryInit)
	now := time.Now().Round(time.Minute)

	queryCmd.PersistentFlags().StringVarP(&requestcontext.Cluster, "cluster", "c", "", "Confluent Cloud Cluster ID")
	queryCmd.MarkPersistentFlagRequired("cluster")

	queryCmd.PersistentFlags().StringVar(&requestcontext.StartTime, "start", now.Add(time.Duration(-1)*time.Hour).Format(ccloudmetrics.TimeFormatStr), "Start Time in the format of "+ccloudmetrics.TimeFormatStr)
	queryCmd.PersistentFlags().StringVar(&requestcontext.EndTime, "end", now.Format(ccloudmetrics.TimeFormatStr), "End Time in the format of "+ccloudmetrics.TimeFormatStr)
	queryCmd.PersistentFlags().IntVar(&requestcontext.LastXmin, "last", 0, "Instead of using start/end time. Query for the last X mins")

	queryCmd.PersistentFlags().StringVar(&requestcontext.Granularity, "gran", string(ccloudmetrics.GranularityOneHour), "Granularity of Metrics. Options are: "+strings.Join(ccloudmetrics.AvailableGranularities, ", "))

	rootCmd.AddCommand(queryCmd)
}

func queryInit() {
}
