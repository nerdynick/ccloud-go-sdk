package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	"github.com/spf13/cobra"
)

const (
	OutputPlain OutputFormat = "plain"
	OutputJSON  OutputFormat = "json"
	OutputCSV   OutputFormat = "csv"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccloudmetrics",
		Short: "Confluent Cloud Metrics API CLI",
	}

	//All Supported output formats
	AvailableOutputFormats = map[string]OutputFormat{
		"plain": OutputPlain,
		"json":  OutputJSON,
		"csv":   OutputCSV,
	}
)

//Global Vars
var (
	verbose           bool
	extraVerbose      bool
	extraExtraVerbose bool
	strOutputFormat   string                    = string(OutputPlain)
	apiContext        ccloudmetrics.APIContext  = ccloudmetrics.NewAPIContext("", "")
	httpContext       ccloudmetrics.HTTPContext = ccloudmetrics.NewHTTPContext()
	requestcontext    RequestContext
)

//Output Formats
type OutputFormat string

func init() {
	cobra.OnInitialize(rootInit)

	log.WithFields(log.Fields{
		"APIContext":  apiContext,
		"HTTPContext": httpContext,
	}).Trace("Initial Contexts")

	//Root Commands
	rootCmd.PersistentFlags().BoolVar(&verbose, "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraVerbose, "vv", false, "Extra Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraExtraVerbose, "vvv", false, "Extra Extra Verbose output")
	rootCmd.PersistentFlags().StringVarP(&apiContext.APIKey, "apikey", "k", "", "API Key")
	rootCmd.MarkPersistentFlagRequired("apikey")
	rootCmd.PersistentFlags().StringVarP(&apiContext.APISecret, "apisecret", "s", "", "API Secret")
	rootCmd.MarkPersistentFlagRequired("apisecret")
	rootCmd.PersistentFlags().StringVarP(&strOutputFormat, "output", "o", strOutputFormat, "Output Format - Available Options: plain, csv, json")
	rootCmd.PersistentFlags().StringVarP(&apiContext.BaseURL, "baseurl", "b", ccloudmetrics.DefaultBaseURL, "API Base Url")
	rootCmd.PersistentFlags().IntVarP(&httpContext.RequestTimeout, "timeout", "t", ccloudmetrics.DefaultRequestTimeout, "HTTP Request Timeout")
	rootCmd.PersistentFlags().StringVarP(&httpContext.UserAgent, "agent", "a", "ccloud-metrics-sdk/go-cli", "HTTP User Agent")
}

func rootInit() {
	requestcontext.OutputFormat = AvailableOutputFormats[strings.ToLower(strOutputFormat)]

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
	req(*cobra.Command, []string, RequestContext, ccloudmetrics.MetricsClient) (bool, error)
	outputPlain() error
	outputJSON(*json.Encoder) error
	outputCSV(*csv.Writer) error
}

func runE(run runFunc) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := getClient()
		results, err := run.req(cmd, args, requestcontext, client)
		if err != nil {
			log.Panic(fmt.Sprintf("Failed to get full results. Error: %s", err.Error()))
			return err
		}

		outputErrs := []error{}
		if err != nil {
			outputErrs = append(outputErrs, err)
		}

		if results {
			switch requestcontext.OutputFormat {
			case OutputCSV:
				writer := csv.NewWriter(os.Stdout)
				defer writer.Flush()
				err := run.outputCSV(writer)
				if err != nil {
					outputErrs = append(outputErrs, err)
				}
				break
			case OutputJSON:
				encoder := json.NewEncoder(os.Stdout)
				err := run.outputJSON(encoder)
				if err != nil {
					outputErrs = append(outputErrs, err)
				}
				break
			case OutputPlain:
				err := run.outputPlain()
				if err != nil {
					outputErrs = append(outputErrs, err)
				}
				break
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
}

func Execute() error {
	return rootCmd.Execute()
}
