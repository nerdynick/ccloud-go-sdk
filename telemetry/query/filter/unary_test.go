package filter

import (
	"testing"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/stretchr/testify/assert"
)

func TestUnaryFilter_And(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := Not(fil1)

	assert.Equal(OpNot, fil.Op, "Filter Op Doesn't Match")
	assert.Equal(OpEq, fil.SubFilter.(FieldFilter).Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.SubFilter.(FieldFilter).Field, "Filter Field doesn't match expected")

	filAnd := fil.And(fil2)

	assert.Equal(OpAnd, filAnd.Op, "Filter Op Doesn't Match")
	assert.Len(filAnd.Filters, 2, "Not all filters where in collection")
}

func TestUnaryFilter_Or(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := Not(fil1)

	assert.Equal(OpNot, fil.Op, "Filter Op Doesn't Match")
	assert.Equal(OpEq, fil.SubFilter.(FieldFilter).Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.SubFilter.(FieldFilter).Field, "Filter Field doesn't match expected")

	filAnd := fil.Or(fil2)

	assert.Equal(OpOr, filAnd.Op, "Filter Op Doesn't Match")
	assert.Len(filAnd.Filters, 2, "Not all filters where in collection")
}
