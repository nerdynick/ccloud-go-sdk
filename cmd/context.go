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
	StartTime         string
	EndTime           string
	Topic             string
	Topics            []string
	BlacklistedTopics []string
	IncludePartitions bool
	Granularity       string
	LastXmin          int
	OutputFormat      OutputFormat
	Metric            string
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

func (r *RequestContext) getGranularity() ccloudmetrics.Granularity {
	g := ccloudmetrics.Granularity(r.Granularity)
	if !g.IsValid() {
		log.Panic(fmt.Sprintf("Granularity is invalid. Was given the value of %s expecting on of %s", r.Granularity, strings.Join(ccloudmetrics.AvailableGranularities, ", ")))
	}
	return g
}
