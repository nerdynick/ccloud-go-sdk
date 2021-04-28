package interval

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/granularity"
	"github.com/rickb777/date/period"
	"github.com/rickb777/date/timespan"
)

type Interval struct {
	timespan.TimeSpan
	withDuration bool
}

func (i Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i Interval) String() string {
	return i.TimeSpan.Format(time.RFC3339, "/", i.withDuration)
}

func (i Interval) MinGranularity() granularity.Granularity {
	for _, g := range granularity.AvailableGranularities {
		if g.IsValidInterval(i) {
			return g
		}
	}
	return granularity.OneMin
}

func Between(start, end time.Time) Interval {
	return Interval{
		TimeSpan:     timespan.NewTimeSpan(start.Round(time.Minute), end.Round(time.Minute)),
		withDuration: false,
	}
}

func StartingFrom(start time.Time, duration time.Duration) Interval {
	return Interval{
		TimeSpan:     timespan.TimeSpanOf(start.Round(time.Minute), duration),
		withDuration: true,
	}
}

func EndingAt(duration time.Duration, end time.Time) Interval {
	return Interval{
		TimeSpan:     timespan.TimeSpanOf(end.Round(time.Minute), -duration),
		withDuration: true,
	}
}

func Of(intervals ...Interval) []Interval {
	return intervals
}

func Parse(value string) (Interval, error) {
	slash := strings.IndexByte(value, '/')
	if slash < 0 {
		return Interval{}, fmt.Errorf("cannot parse %q because there is no separator '/'", value)
	}

	start := value[:slash]
	rest := value[slash+1:]

	if rest == "" {
		return Interval{}, fmt.Errorf("cannot parse %q because there is end time or duration", value)
	}

	if start[0] == 'P' {
		p, err := period.Parse(start)
		if err != nil {
			return Interval{}, err
		}
		t, err := time.Parse(time.RFC3339, rest)
		if err != nil {
			return Interval{}, err
		}
		return EndingAt(p.DurationApprox(), t), nil
	} else if rest[0] == 'P' {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return Interval{}, err
		}
		p, err := period.Parse(rest)
		if err != nil {
			return Interval{}, err
		}
		return StartingFrom(t, p.DurationApprox()), nil
	} else {
		s, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return Interval{}, err
		}

		e, err := time.Parse(time.RFC3339, rest)
		if err != nil {
			return Interval{}, err
		}

		return Between(s, e), nil
	}

}
