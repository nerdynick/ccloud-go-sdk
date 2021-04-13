package filter

import "github.com/nerdynick/ccloud-go-sdk/telemetry/labels"

const (
	//OpNot is a static def for NOT Operand
	OpNot string = "NOT"
)

type UnaryFilter struct {
	Op        string `json:"op"`
	SubFilter Filter `json:"filter"`
}

func (fil UnaryFilter) And(filters ...Filter) CompoundFilter {
	return And(fil).Add(filters...)
}
func (fil UnaryFilter) AndEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(EqualTo(field, value))
}
func (fil UnaryFilter) AndNotEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(NotEqualTo(field, value))
}
func (fil UnaryFilter) AndGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.And(GreaterThan(field, value))
}
func (fil UnaryFilter) AndNotGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.And(NotGreaterThan(field, value))
}
func (fil UnaryFilter) AndGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(GreaterThanOrEqualTo(field, value))
}
func (fil UnaryFilter) AndNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(NotGreaterThanOrEqualTo(field, value))
}

func (fil UnaryFilter) Or(filters ...Filter) CompoundFilter {
	return Or(fil).Add(filters...)
}
func (fil UnaryFilter) OrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(EqualTo(field, value))
}
func (fil UnaryFilter) OrNotEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotEqualTo(field, value))
}
func (fil UnaryFilter) OrGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.Or(GreaterThan(field, value))
}
func (fil UnaryFilter) OrNotGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotGreaterThan(field, value))
}
func (fil UnaryFilter) OrGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(GreaterThanOrEqualTo(field, value))
}
func (fil UnaryFilter) OrNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotGreaterThanOrEqualTo(field, value))
}
