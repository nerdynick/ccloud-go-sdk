package filter

import (
	"testing"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/stretchr/testify/assert"
)

func TestEqualTo(t *testing.T) {
	assert := assert.New(t)

	fil := EqualTo(labels.MetricTopic, "testing")

	assert.Equal(OpEq, fil.Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.Field, "Filter Field doesn't match expected")
	assert.Equal("testing", fil.Value, "Filter Field doesn't match expected")
}

func TestInterface(t *testing.T) {
	assert := assert.New(t)
	var fil Filter

	fil = FieldFilter{}
	assert.NotNil(fil)
	fil = CompoundFilter{}
	assert.NotNil(fil)
	fil = CompoundFilter{}
	assert.NotNil(fil)
}

func TestNotEqualTo(t *testing.T) {
	assert := assert.New(t)

	fil := NotEqualTo(labels.MetricTopic, "testing")

	assert.Equal(OpNot, fil.Op, "Filter Op Doesn't Match")
	assert.Equal(OpEq, fil.SubFilter.(FieldFilter).Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.SubFilter.(FieldFilter).Field, "Filter Field doesn't match expected")
	assert.Equal("testing", fil.SubFilter.(FieldFilter).Value, "Filter Field doesn't match expected")
}

func TestGreaterThan(t *testing.T) {
	assert := assert.New(t)

	fil := GreaterThan(labels.MetricTopic, "testing")

	assert.Equal(OpGt, fil.Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.Field, "Filter Field doesn't match expected")
	assert.Equal("testing", fil.Value, "Filter Field doesn't match expected")
}

func TestNotGreaterThan(t *testing.T) {
	assert := assert.New(t)

	fil := NotGreaterThan(labels.MetricTopic, "testing")

	assert.Equal(OpNot, fil.Op, "Filter Op Doesn't Match")
	assert.Equal(OpGt, fil.SubFilter.(FieldFilter).Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.SubFilter.(FieldFilter).Field, "Filter Field doesn't match expected")
	assert.Equal("testing", fil.SubFilter.(FieldFilter).Value, "Filter Field doesn't match expected")
}

func TestGreaterThanOrEqualTo(t *testing.T) {
	assert := assert.New(t)

	fil := GreaterThanOrEqualTo(labels.MetricTopic, "testing")

	assert.Equal(OpGte, fil.Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.Field, "Filter Field doesn't match expected")
	assert.Equal("testing", fil.Value, "Filter Field doesn't match expected")
}

func TestNotGreaterThanOrEqualTo(t *testing.T) {
	assert := assert.New(t)

	fil := NotGreaterThanOrEqualTo(labels.MetricTopic, "testing")

	assert.Equal(OpNot, fil.Op, "Filter Op Doesn't Match")
	assert.Equal(OpGte, fil.SubFilter.(FieldFilter).Op, "Filter OP doesn't match expected")
	assert.Equal(labels.MetricTopic, fil.SubFilter.(FieldFilter).Field, "Filter Field doesn't match expected")
	assert.Equal("testing", fil.SubFilter.(FieldFilter).Value, "Filter Field doesn't match expected")
}

func TestAnd(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := And(fil1, fil2)

	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}

func TestAllOf(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := AllOf(fil1, fil2)

	assert.Equal(OpAnd, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}

func TestOr(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := Or(fil1, fil2)

	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}

func TestOneOf(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := OneOf(fil1, fil2)

	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}

func TestAnyOf(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := AnyOf(fil1, fil2)

	assert.Equal(OpOr, fil.Op, "Filter Op Doesn't Match")
	assert.Len(fil.Filters, 2, "Not all filters where in collection")
}

func TestNotAnyOf(t *testing.T) {
	assert := assert.New(t)

	fil1 := EqualTo(labels.MetricTopic, "testing1")
	fil2 := EqualTo(labels.MetricTopic, "testing2")
	fil := NotAnyOf(fil1, fil2)

	assert.Equal(OpNot, fil.Op, "Filter Op Doesn't Match")
	assert.Equal(OpOr, fil.SubFilter.(CompoundFilter).Op, "Filter Op Doesn't Match")
	assert.Len(fil.SubFilter.(CompoundFilter).Filters, 2, "Not all filters where in collection")
}
