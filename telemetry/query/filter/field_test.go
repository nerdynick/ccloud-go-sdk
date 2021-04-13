package filter

import (
	"testing"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/stretchr/testify/assert"
)

func TestFieldFilter_And(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")

	assert.Equal(OpEq, fil1.Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil1.Field, "Filter Field doesn't match expected")

	filAnd := fil1.And(fil2)

	assert.Equal(OpAnd, filAnd.Op, "Filter Op Doesn't Match")
	assert.Len(filAnd.Filters, 2, "Not all filters where in collection")
}

func TestFieldFilter_Or(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")

	assert.Equal(OpEq, fil1.Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil1.Field, "Filter Field doesn't match expected")

	filAnd := fil1.Or(fil2)

	assert.Equal(OpOr, filAnd.Op, "Filter Op Doesn't Match")
	assert.Len(filAnd.Filters, 2, "Not all filters where in collection")
}

func TestFieldFilter_Not(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")

	assert.Equal(OpEq, fil1.Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil1.Field, "Filter Field doesn't match expected")

	filAnd := fil1.Not()

	assert.Equal(OpNot, filAnd.Op, "Filter Op Doesn't Match")
	assert.Equal(OpEq, filAnd.SubFilter.(FieldFilter).Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, filAnd.SubFilter.(FieldFilter).Field, "Filter Field doesn't match expected")
}
