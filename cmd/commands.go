package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccloudmetrics",
		Short: "Confluent Cloud Metrics API CLI",
	}
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get metrics back from the API",
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List information from the API",
	}

	availableMetrics = &cobra.Command{
		Use:   "metrics",
		Short: "List currently available metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := ccloudmetrics.NewClientFromContext(&apiContext, &httpContext)

			res, err := client.GetAvailableMetrics()
			if err != nil {
				return err
			}
			writer := csv.NewWriter(os.Stdout)
			for _, metric := range res {
				labels := []string{}
				for _, label := range metric.Labels {
					labels = append(labels, label.Name)
				}

				if !outputJSON && !outputCSV {
					fmt.Printf("Name:      %s\n", metric.Name)
					fmt.Printf("Desc:      %s\n", metric.Desc)
					fmt.Printf("Type:      %s\n", metric.Type)
					fmt.Printf("LifeCycle: %s\n", metric.LifecycleStage)
					fmt.Printf("Labels:    %s\n", strings.Join(labels, ","))
					fmt.Println()
				}

				if outputJSON {
					j, err := json.Marshal(metric)
					if err != nil {
						return err
					}
					fmt.Printf("%s", j)
					fmt.Println()
				}
				if outputCSV {
					writer.Write([]string{
						metric.Name,
						metric.Desc,
						metric.Type,
						metric.LifecycleStage,
						strings.Join(labels, ";"),
					})
				}
			}
			writer.Flush()

			return nil
		},
	}
)

var (
	verbose     bool
	outputJSON  bool
	outputCSV   bool
	apiContext  ccloudmetrics.APIContext
	httpContext ccloudmetrics.HTTPContext
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "JSON output")
	rootCmd.PersistentFlags().BoolVar(&outputCSV, "csv", false, "CSV output")

	rootCmd.PersistentFlags().StringVarP(&apiContext.APIKey, "apikey", "k", "", "API Key")
	rootCmd.MarkPersistentFlagRequired("apikey")

	rootCmd.PersistentFlags().StringVarP(&apiContext.APISecret, "apisecret", "s", "", "API Secret")
	rootCmd.MarkPersistentFlagRequired("apisecret")

	rootCmd.PersistentFlags().StringVarP(&apiContext.BaseURL, "baseurl", "b", ccloudmetrics.DefaultBaseURL, "API Base Url")

	rootCmd.PersistentFlags().IntVarP(&httpContext.RequestTimeout, "timeout", "t", ccloudmetrics.DefaultRequestTimeout, "HTTP Request Timeout")

	rootCmd.PersistentFlags().StringVarP(&httpContext.UserAgent, "agent", "a", ccloudmetrics.DefaultUserAgent, "HTTP User Agent")

	listCmd.AddCommand(availableMetrics)

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
