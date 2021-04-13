package filter

import (
	"testing"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/stretchr/testify/assert"
)

func TestCompoundFilter_Add(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil3 := EqualTo(labels.MetricTopic, "testing3")
	fil := And(fil1, fil2)

	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")

	fil = fil.Add(fil3)

	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 3, "Not all filters where in collection")
}

func TestCompoundFilter_And(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil3 := EqualTo(labels.MetricTopic, "testing3")

	//Test 1
	fil := And(fil1, fil2)
	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")

	fil = fil.And(fil3)
	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 3, "Not all filters where in collection")

	//Test 2
	fil = Or(fil1, fil2)
	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")

	fil = fil.And(fil3)
	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}

func TestCompoundFilter_Or(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil3 := EqualTo(labels.MetricTopic, "testing3")

	//Test 1
	fil := Or(fil1, fil2)
	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")

	fil = fil.Or(fil3)
	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 3, "Not all filters where in collection")

	//Test 2
	fil = And(fil1, fil2)
	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")

	fil = fil.Or(fil3)
	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}
