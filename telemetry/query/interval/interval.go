package interval

import (
	"encoding/json"
	"time"

	"github.com/rickb777/date/timespan"
)

type Interval struct {
	timespan.TimeSpan
	withDuration bool
}

func (i *Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i Interval) String() string {
	return i.TimeSpan.Format(time.RFC3339, "/", i.withDuration)
}

func Between(start, end time.Time) Interval {
	return Interval{
		TimeSpan:     timespan.NewTimeSpan(start, end),
		withDuration: false,
	}
}

func StartingFrom(start time.Time, duration time.Duration) Interval {
	return Interval{
		TimeSpan:     timespan.TimeSpanOf(start, duration),
		withDuration: true,
	}
}

func EndingAt(duration time.Duration, end time.Time) Interval {
	return Interval{
		TimeSpan:     timespan.TimeSpanOf(end, -duration),
		withDuration: true,
	}
}

func Of(intervals ...Interval) []Interval {
	return intervals
}
