package granularity

import (
	"encoding/json"
	"math"
	"time"

	"github.com/rickb777/date/period"
)

var (
	//OneMin is a static def for the 1 minute granularity string
	OneMin Granularity = newGranularity("PT1M")
	//FiveMin is a static def for the 5 minute granularity string
	FiveMin Granularity = newGranularity("PT5M")
	//FifteenMin is a static def for the 15 minute granularity string
	FifteenMin Granularity = newGranularity("PT15M")
	//ThirtyMin is a static def for the 30 minute granularity string
	ThirtyMin Granularity = newGranularity("PT30M")
	//OneHour is a static def for the 1 hour granularity string
	OneHour Granularity = newGranularity("PT1H")
	//FourHours is a static def for the 4 hours granularity string
	FourHours Granularity = newGranularity("PT4H")
	//SixHours is a static def for the 6 hours granularity string
	SixHours Granularity = newGranularity("PT6H")
	//TwelveHours is a static def for the 12 hours granularity string
	TwelveHours Granularity = newGranularity("PT12H")
	//OneDay is a static def for the 1 day granularity string
	OneDay Granularity = newGranularity("P1D")
	//All is a static def for the ALL granularity string
	All Granularity = newGranularity("ALL")

	//AvailableGranularities is a collection of all available Granularities
	AvailableGranularities []Granularity = []Granularity{
		OneMin,
		FiveMin,
		FifteenMin,
		ThirtyMin,
		OneHour,
		OneDay,
		All,
	}

	maxDuration time.Duration = (math.MaxInt64 * time.Nanosecond)
)

//Granularity string type to extend extra helper functions
type Granularity struct {
	string
	time.Duration
	MaxDuration time.Duration
	Period      period.Period
}

func (g Granularity) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g Granularity) String() string {
	return g.string
}

func (g Granularity) Equals(gran Granularity) bool {
	return g.Duration == gran.Duration
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

func newGranularity(g string) Granularity {
	if g == "ALL" {
		p, _ := period.NewOf(maxDuration)
		return Granularity{
			string:      g,
			Duration:    maxDuration,
			MaxDuration: maxDuration,
			Period:      p,
		}
	}
	p, _ := period.Parse(g)
	d, _ := p.Duration()

	//Find max duration for a given Granularity. This is based off the documented maxs for Grandularities.
	var maxDur time.Duration
	if d <= time.Minute {
		maxDur = time.Hour * 6
	} else if d <= (5 * time.Minute) {
		maxDur = time.Hour * 24
	} else if d <= (15 * time.Minute) {
		maxDur = time.Hour * (24 * 4)
	} else if d <= (30 * time.Minute) {
		maxDur = time.Hour * (24 * 7)
	} else {
		maxDur = maxDuration
	}

	return Granularity{
		string:      g,
		Duration:    d,
		MaxDuration: maxDur,
		Period:      p,
	}
}
