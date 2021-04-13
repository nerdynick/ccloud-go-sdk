package granularity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStaticDefs_MaxDuration(t *testing.T) {
	assert := assert.New(t)

	assert.Equal((6 * time.Hour).Nanoseconds(), OneMin.MaxDuration.Nanoseconds(), "OneMin has wrong MaxDur")
	assert.Equal((24 * time.Hour).Nanoseconds(), FiveMin.MaxDuration.Nanoseconds(), "FiveMin has wrong MaxDur")
	assert.Equal(4*(24*time.Hour).Nanoseconds(), FifteenMin.MaxDuration.Nanoseconds(), "FifteenMin has wrong MaxDur")
	assert.Equal(7*(24*time.Hour).Nanoseconds(), ThirtyMin.MaxDuration.Nanoseconds(), "ThirtyMin has wrong MaxDur")
	assert.Equal(maxDuration.Nanoseconds(), OneHour.MaxDuration.Nanoseconds(), "OneHour has wrong MaxDur")
	assert.Equal(maxDuration.Nanoseconds(), FourHours.MaxDuration.Nanoseconds(), "FourHours has wrong MaxDur")
	assert.Equal(maxDuration.Nanoseconds(), SixHours.MaxDuration.Nanoseconds(), "SixHours has wrong MaxDur")
	assert.Equal(maxDuration.Nanoseconds(), TwelveHours.MaxDuration.Nanoseconds(), "TwelveHours has wrong MaxDur")
	assert.Equal(maxDuration.Nanoseconds(), OneDay.MaxDuration.Nanoseconds(), "OneDay has wrong MaxDur")

}

func TestStaticDefs_Duration(t *testing.T) {
	assert := assert.New(t)

	assert.Equal((1 * time.Minute), OneMin.Duration, "OneMin has wrong Duration")
	assert.Equal((5 * time.Minute), FiveMin.Duration, "FiveMin has wrong Duration")
	assert.Equal((15 * time.Minute), FifteenMin.Duration, "FifteenMin has wrong Duration")
	assert.Equal((30 * time.Minute), ThirtyMin.Duration, "ThirtyMin has wrong Duration")
	assert.Equal((1 * time.Hour), OneHour.Duration, "OneHour has wrong Duration")
	assert.Equal((4 * time.Hour), FourHours.Duration, "FourHours has wrong Duration")
	assert.Equal((6 * time.Hour), SixHours.Duration, "SixHours has wrong Duration")
	assert.Equal((12 * time.Hour), TwelveHours.Duration, "TwelveHours has wrong Duration")
	assert.Equal((24 * time.Hour), OneDay.Duration, "OneDay has wrong Duration")
}
