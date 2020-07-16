package cmd

import (
	"encoding/csv"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/pflag"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccloudmetrics",
		Short: "Confluent Cloud Metrics API CLI",
	}
)

//Global Vars
var (
	verbose      bool
	extraVerbose bool
	outputJSON   bool
	outputCSV    bool
	apiContext   ccloudmetrics.APIContext
	httpContext  ccloudmetrics.HTTPContext
)

//Common Command Vars
var (
	cluster   string
	metric    string
	startTime string
	endTime   string
)

func init() {
	cobra.OnInitialize(onInit)
	//Root Commands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraVerbose, "vvvv", false, "Extra Verbose output")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "JSON output")
	rootCmd.PersistentFlags().BoolVar(&outputCSV, "csv", false, "CSV output")
	rootCmd.PersistentFlags().StringVarP(&apiContext.APIKey, "apikey", "k", "", "API Key")
	rootCmd.MarkPersistentFlagRequired("apikey")
	rootCmd.PersistentFlags().StringVarP(&apiContext.APISecret, "apisecret", "s", "", "API Secret")
	rootCmd.MarkPersistentFlagRequired("apisecret")
	rootCmd.PersistentFlags().StringVarP(&apiContext.BaseURL, "baseurl", "b", ccloudmetrics.DefaultBaseURL, "API Base Url")
	rootCmd.PersistentFlags().IntVarP(&httpContext.RequestTimeout, "timeout", "t", ccloudmetrics.DefaultRequestTimeout, "HTTP Request Timeout")
	rootCmd.PersistentFlags().StringVarP(&httpContext.UserAgent, "agent", "a", "ccloud-metrics-sdk/go-cli", "HTTP User Agent")
}

func onInit() {
	if verbose || extraVerbose {
		log.SetLevel(log.InfoLevel)

		if extraVerbose {
			log.SetReportCaller(true)
			log.SetLevel(log.TraceLevel)
		}
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func addClusterFlag(flagSet *pflag.FlagSet) {
	topicsForMetric.Flags().StringVarP(&cluster, "cluster", "c", "", "Confluent Cloud Cluster ID")
	topicsForMetric.MarkFlagRequired("cluster")
}

func addMetricFlag(flagSet *pflag.FlagSet) {
	topicsForMetric.Flags().StringVarP(&metric, "metric", "m", "", "Metric to fetch available topics for")
	topicsForMetric.MarkFlagRequired("metric")
}

func getClient() ccloudmetrics.MetricsClient {
	return ccloudmetrics.NewClientFromContext(&apiContext, &httpContext)
}

type runFunc interface {
	req(*cobra.Command, []string, ccloudmetrics.MetricsClient) error
	outputPlain() error
	outputJSON(*json.Encoder) error
	outputCSV(*csv.Writer) error
}

func runE(run runFunc) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := getClient()
		err := run.req(cmd, args, client)

		if err != nil {
			return err
		}
		if outputCSV {
			writer := csv.NewWriter(os.Stdout)
			err := run.outputCSV(writer)
			if err != nil {
				return err
			}
			writer.Flush()
		}
		if outputJSON {
			encoder := json.NewEncoder(os.Stdout)
			return run.outputJSON(encoder)

		}
		if !outputCSV && !outputJSON {
			return run.outputPlain()
		}
		return nil
	}
}

func Execute() error {
	return rootCmd.Execute()
}
