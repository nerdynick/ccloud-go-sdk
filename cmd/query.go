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

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query metrics back from the API",
}

type Query struct {
	Results []ccloudmetrics.QueryData
	request func(*cobra.Command, []string, ccloudmetrics.MetricsClient, time.Time, time.Time) ([]ccloudmetrics.QueryData, error)
}

func (q *Query) req(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient) error {
	sTime, err := time.Parse(ccloudmetrics.TimeFormatStr, startTime)
	if err != nil {
		return nil
	}
	eTime, err := time.Parse(ccloudmetrics.TimeFormatStr, endTime)
	if err != nil {
		return nil
	}

	res, err := q.request(cmd, args, client, sTime, eTime)
	q.Results = res
	return err
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

var (
	granularity string
	lastXmin    int
)

func init() {
	cobra.OnInitialize(queryOnInit)
	queryCmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", "", "Confluent Cloud Cluster ID")
	queryCmd.MarkPersistentFlagRequired("cluster")

	queryCmd.PersistentFlags().StringVar(&startTime, "start", time.Now().Add(time.Duration(-1)*time.Hour).Format(ccloudmetrics.TimeFormatStr), "Start Time in the format of "+ccloudmetrics.TimeFormatStr)
	queryCmd.PersistentFlags().StringVar(&endTime, "end", time.Now().Format(ccloudmetrics.TimeFormatStr), "End Time in the format of "+ccloudmetrics.TimeFormatStr)
	queryCmd.PersistentFlags().IntVar(&lastXmin, "last", 0, "Instead of using start/end time. Query for the last X mins")

	queryCmd.PersistentFlags().StringVar(&granularity, "gran", ccloudmetrics.GranularityOneHour, "Granularity of Metrics. Options are: "+strings.Join(ccloudmetrics.AvailableGranularities, ", "))

	rootCmd.AddCommand(queryCmd)
}

func queryOnInit() {
	if lastXmin > 0 {
		startTime = time.Now().Add(time.Duration(-lastXmin) * time.Minute).Format(ccloudmetrics.TimeFormatStr)
	}
}
