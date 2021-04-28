package interval

import (
	"testing"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/granularity"
	"github.com/stretchr/testify/assert"
)

func TestParseRange(t *testing.T) {
	assert := assert.New(t)

	interval := "2021-04-19T15:15:00-06:00/2021-04-20T16:15:00-06:00"
	i, e := Parse(interval)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(interval, i.String())
}

func TestParseStart(t *testing.T) {
	assert := assert.New(t)

	interval := "2021-04-19T15:15:00-06:00/PT1H"
	i, e := Parse(interval)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(interval, i.String())
}

func TestParseEnd(t *testing.T) {
	assert := assert.New(t)

	interval := "PT1H/2021-04-20T16:15:00-06:00"
	i, e := Parse(interval)
	if e != nil {
		t.Error(e)
	}

	assert.Equal("2021-04-20T15:15:00-06:00/PT1H", i.String())
}

func TestParseRangeWithError(t *testing.T) {
	assert := assert.New(t)

	_, e := Parse("2021-04-19T15:15/2021-04-20T16:15:00-06:00")
	assert.Error(e)

	_, e = Parse("2021-04-19T15:15:00-06:00/2021-04-20T16:15")
	assert.Error(e)
}

func TestParseStartWithError(t *testing.T) {
	assert := assert.New(t)

	_, e := Parse("2021-04-19T15:15:00-06:00/P")
	assert.Error(e)

	_, e = Parse("2021-04-19T15:15/PT1H")
	assert.Error(e)
}

func TestParseEndWithError(t *testing.T) {
	assert := assert.New(t)

	_, e := Parse("P/2021-04-20T16:15:27-06:00")
	assert.Error(e)

	_, e = Parse("PT1H/2021-04-20T16:15")
	assert.Error(e)
}

func TestParseOnlyHalf(t *testing.T) {
	assert := assert.New(t)

	_, e := Parse("P1H/")
	assert.Error(e)
}
func TestParseNoSlash(t *testing.T) {
	assert := assert.New(t)

	_, e := Parse("P")
	assert.Error(e)
}

func TestMinGranularity(t *testing.T) {
	assert := assert.New(t)

	interval, e := Parse("2021-04-19T15:15:00-06:00/PT1H")
	if e != nil {
		t.Error(e)
	}

	gran := interval.MinGranularity()
	assert.Equal(granularity.OneMin, gran)

}
func TestValidGranularity(t *testing.T) {
	assert := assert.New(t)

	interval, e := Parse("2021-04-19T15:15:00-06:00/P8DT")
	if e != nil {
		t.Error(e)
	}

	assert.False(interval.IsValidGranularity(granularity.OneMin))
	assert.False(interval.IsValidGranularity(granularity.FiveMin))
	assert.False(interval.IsValidGranularity(granularity.FifteenMin))
	assert.False(interval.IsValidGranularity(granularity.ThirtyMin))
	assert.True(interval.IsValidGranularity(granularity.OneHour))

}
