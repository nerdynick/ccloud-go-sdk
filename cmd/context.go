package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
)

type RequestContext struct {
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
	OutputFormat      OutputFormat
}

func (r *RequestContext) getStartTime() time.Time {
	if r.LastXmin > 0 {
		return time.Now().Add(time.Duration(-r.LastXmin) * time.Minute)
	}
	res, err := time.Parse(ccloudmetrics.TimeFormatStr, r.StartTime)
	if err != nil {
		log.Panic(fmt.Sprintf("Start Time is invalid. Times must be provided in the %s format. Was given %s", ccloudmetrics.TimeFormatStr, r.StartTime))
	}
	return res
}
func (r *RequestContext) getEndTime() time.Time {
	res, err := time.Parse(ccloudmetrics.TimeFormatStr, r.EndTime)
	if err != nil {
		log.Panic(fmt.Sprintf("End Time is invalid. Times must be provided in the %s format. Was given %s", ccloudmetrics.TimeFormatStr, r.EndTime))
	}
	return res
}
func (r *RequestContext) getMetric() ccloudmetrics.Metric {
	metrics, err := getClient().GetAvailableMetrics()
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to get all Available Metrics. Got error %s", err.Error()))
	}
	metricNames := []string{}

	for _, metric := range metrics {
		metricNames = append(metricNames, metric.Name)
		if metric.Matches(r.Metric) {
			return metric
		}
	}

	log.Panic(fmt.Sprintf("Metric is invalid. Got %s but only have available %s", r.Metric, strings.Join(metricNames, ", ")))
	return ccloudmetrics.Metric{}
}
func (r *RequestContext) getMetrics() []ccloudmetrics.Metric {
	metrics, err := getClient().GetAvailableMetrics()
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to get all Available Metrics. Got error %s", err.Error()))
	}
	validMetrics := []ccloudmetrics.Metric{}

	for _, metric := range metrics {
		for _, m := range r.Metrics {
			if metric.Matches(m) {
				validMetrics = append(validMetrics, metric)
			}
		}
	}
	return validMetrics
}
func (r *RequestContext) getGranularity() ccloudmetrics.Granularity {
	g := ccloudmetrics.Granularity(r.Granularity)
	if !g.IsValid() {
		log.Panic(fmt.Sprintf("Granularity is invalid. Was given the value of %s expecting on of %s", r.Granularity, strings.Join(ccloudmetrics.AvailableGranularities, ", ")))
	}
	return g
}
