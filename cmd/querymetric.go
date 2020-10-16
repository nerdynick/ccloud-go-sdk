package cmd

import (
	"fmt"
	"strings"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	validMetrics := []string{}
	for _, m := range ccloudmetrics.KnownMetrics {
		validMetrics = append(validMetrics, fmt.Sprintf("%-60v  -or-  %10v", m.Name, m.ShortName()))
	}

	queryCmd.AddCommand(
		&cobra.Command{
			Use:     "metric",
			Aliases: []string{"metrics"},
			Short:   "Query results for a given metric. \nAvailable Known Metrics: \n\t" + strings.Join(validMetrics, "\n\t"),
			Args:    cobra.MinimumNArgs(1),
			RunE: runE(&Query{
				request: func(cmd *cobra.Command, args []string, client ccloudmetrics.MetricsClient, context RequestContext) ([]ccloudmetrics.QueryData, error) {
					metrics := getMetrics(client, args...)
					log.WithField("metrics", metrics).Info("Metrics")
					if len(metrics) > 1 {
						return client.QueryMetrics(context.Cluster, metrics, context.getGranularity(), context.getStartTime(), context.getEndTime())
					} else if len(metrics) == 1 {
						return client.QueryMetric(context.Cluster, metrics[0], context.getGranularity(), context.getStartTime(), context.getEndTime())
					} else {
						return nil, nil
					}
				},
			}),
		})
}
