package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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
	verbose           bool
	extraVerbose      bool
	extraExtraVerbose bool
	outputJSON        bool
	outputCSV         bool
	apiContext        ccloudmetrics.APIContext  = ccloudmetrics.NewAPIContext("", "")
	httpContext       ccloudmetrics.HTTPContext = ccloudmetrics.NewHTTPContext()
)

//Common Command Vars
var (
	context requestContext
)

type requestContext struct {
	Cluster           string
	Metric            string
	Metrics           []string
	StartTime         string
	EndTime           string
	Topic             string
	Topics            []string
	BlacklistedTopics []string
	IncludePartitions bool
	Granularity       string
	LastXmin          int
}

func (r requestContext) getStartTime() time.Time {
	if r.LastXmin > 0 {
		return time.Now().Add(time.Duration(-r.LastXmin) * time.Minute)
	}
	res, err := time.Parse(ccloudmetrics.TimeFormatStr, r.StartTime)
	if err != nil {
		log.Panic(fmt.Sprintf("Start Time is invalid. Times must be provided in the %s format. Was given %s", ccloudmetrics.TimeFormatStr, context.StartTime))
	}
	return res
}
func (r requestContext) getEndTime() time.Time {
	res, err := time.Parse(ccloudmetrics.TimeFormatStr, r.EndTime)
	if err != nil {
		log.Panic(fmt.Sprintf("End Time is invalid. Times must be provided in the %s format. Was given %s", ccloudmetrics.TimeFormatStr, context.EndTime))
	}
	return res
}
func (r requestContext) getMetric() ccloudmetrics.Metric {
	metrics, err := getClient().GetAvailableMetrics()
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to get all Available Metrics. Got error %s", err.Error()))
	}
	metricNames := []string{}

	for _, metric := range metrics {
		metricNames = append(metricNames, metric.Name)
		if metric.Name == r.Metric {
			return metric
		}
	}

	log.Panic(fmt.Sprintf("Metric is invalid. Got %s but only have available %s", context.Metric, strings.Join(metricNames, ", ")))
	return ccloudmetrics.Metric{}
}
func (r requestContext) getMetrics() []ccloudmetrics.Metric {
	metrics, err := getClient().GetAvailableMetrics()
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to get all Available Metrics. Got error %s", err.Error()))
	}
	validMetrics := []ccloudmetrics.Metric{}

	for _, metric := range metrics {
		for _, m := range context.Metrics {
			if metric.Name == m {
				validMetrics = append(validMetrics, metric)
			}
		}
	}
	return validMetrics
}
func (r requestContext) getGranularity() ccloudmetrics.Granularity {
	g := ccloudmetrics.Granularity(context.Granularity)
	if !g.IsValid() {
		log.Panic(fmt.Sprintf("Granularity is invalid. Was given the value of %s expecting on of %s", context.Granularity, strings.Join(ccloudmetrics.AvailableGranularities, ", ")))
	}
	return g
}

func init() {
	cobra.OnInitialize(onInit)

	log.WithFields(log.Fields{
		"APIContext":  apiContext,
		"HTTPContext": httpContext,
	}).Trace("Initial Contexts")

	//Root Commands
	rootCmd.PersistentFlags().BoolVar(&verbose, "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraVerbose, "vv", false, "Extra Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraExtraVerbose, "vvv", false, "Extra Extra Verbose output")
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
	if verbose || extraVerbose || extraExtraVerbose {
		log.SetLevel(log.InfoLevel)

		if extraVerbose {
			log.SetLevel(log.DebugLevel)
		}

		if extraExtraVerbose {
			log.SetReportCaller(true)
			log.SetLevel(log.TraceLevel)
		}
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func getClient() ccloudmetrics.MetricsClient {
	return ccloudmetrics.NewClientFromContext(apiContext, httpContext)
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
			log.Panic(fmt.Sprintf("Failed to get full results. Error: %s", err.Error()))
		}
		return err
	}
}

func Execute() error {
	return rootCmd.Execute()
}
