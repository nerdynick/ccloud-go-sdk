package ccloudmetrics

import "time"

const (
	//GranularityOneMin is a static def for the 1 minute grantularity string
	GranularityOneMin Granularity = "PT1M"
	//GranularityFiveMin is a static def for the 5 minute grantularity string
	GranularityFiveMin Granularity = "PT5M"
	//GranularityFifteenMin is a static def for the 15 minute grantularity string
	GranularityFifteenMin Granularity = "PT15M"
	//GranularityThirtyMin is a static def for the 30 minute grantularity string
	GranularityThirtyMin Granularity = "PT30M"
	//GranularityOneHour is a static def for the 1 hour grantularity string
	GranularityOneHour Granularity = "PT1H"
	//GranularityAll is a static def for the ALL grantularity string
	GranularityAll Granularity = "ALL"
)

var (
	//AvailableGranularities is a collection of all available Granularities
	AvailableGranularities []string = []string{
		string(GranularityOneMin),
		string(GranularityFiveMin),
		string(GranularityFifteenMin),
		string(GranularityThirtyMin),
		string(GranularityOneHour),
		string(GranularityAll),
	}
)

//Granularity string type to extend extra helper functions
type Granularity string

func (g Granularity) String() string {
	return string(g)
}

//IsValid checks in the current Granularity is a valid, available, and known Granularity
func (g Granularity) IsValid() bool {
	for _, l := range AvailableGranularities {
		if g.Equals(l) {
			return true
		}
	}
	return false
}

//Equals tests if the Granularity is equal to a string representation
func (g Granularity) Equals(str string) bool {
	if string(g) == str {
		return true
	}
	return false
}

//GetStartTimeFromGranularity is a utility func to get a Start time given a granularity
func (g Granularity) GetStartTimeFromGranularity(t time.Time) time.Time {
	switch g {
	case GranularityOneMin:
		return t.Add(time.Duration(-1) * time.Minute)
	case GranularityFiveMin:
		return t.Add(time.Duration(-5) * time.Minute)
	case GranularityFifteenMin:
		return t.Add(time.Duration(-15) * time.Minute)
	case GranularityThirtyMin:
		return t.Add(time.Duration(-30) * time.Minute)
	case GranularityOneHour:
		return t.Add(time.Duration(-1) * time.Hour)
	}
	return time.Now()
}
