package filter

import "github.com/nerdynick/ccloud-go-sdk/telemetry/labels"

const (
	//OpEq is a static def for EQ Operand
	OpEq  string = "EQ"
	OpGt  string = "GT"
	OpGte string = "GTE"
)

// Filter structure
type FieldFilter struct {
	Op    string       `json:"op"`
	Field labels.Label `json:"field,omitempty"`
	Value string       `json:"value"`
}

func (fil FieldFilter) Not() UnaryFilter {
	return Not(fil)
}

func (fil FieldFilter) And(filters ...Filter) CompoundFilter {
	return And(fil).Add(filters...)
}
func (fil FieldFilter) AndEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(EqualTo(field, value))
}
func (fil FieldFilter) AndNotEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(NotEqualTo(field, value))
}
func (fil FieldFilter) AndGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.And(GreaterThan(field, value))
}
func (fil FieldFilter) AndNotGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.And(NotGreaterThan(field, value))
}
func (fil FieldFilter) AndGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(GreaterThanOrEqualTo(field, value))
}
func (fil FieldFilter) AndNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(NotGreaterThanOrEqualTo(field, value))
}

func (fil FieldFilter) Or(filters ...Filter) CompoundFilter {
	return Or(fil).Add(filters...)
}
func (fil FieldFilter) OrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(EqualTo(field, value))
}
func (fil FieldFilter) OrNotEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotEqualTo(field, value))
}
func (fil FieldFilter) OrGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.Or(GreaterThan(field, value))
}
func (fil FieldFilter) OrNotGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotGreaterThan(field, value))
}
func (fil FieldFilter) OrGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(GreaterThanOrEqualTo(field, value))
}
func (fil FieldFilter) OrNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotGreaterThanOrEqualTo(field, value))
}
