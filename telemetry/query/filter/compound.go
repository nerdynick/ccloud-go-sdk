package filter

import "github.com/nerdynick/ccloud-go-sdk/telemetry/labels"

const (
	//OpAnd is a static def for AND Operand
	OpAnd string = "AND"
	//OpOr is a static def for OR Operand
	OpOr string = "OR"
)

// CompoundFilter to use for a query
type CompoundFilter struct {
	Op      string   `json:"op"`
	Filters []Filter `json:"filters"`
}

func (fil CompoundFilter) Not() UnaryFilter {
	return Not(fil)
}
func (fil CompoundFilter) Add(filters ...Filter) CompoundFilter {
	fil.Filters = append(fil.Filters, filters...)
	return fil
}

func (fil CompoundFilter) And(filters ...Filter) CompoundFilter {
	if fil.Op == OpAnd {
		return fil.Add(filters...)
	} else {
		return And(fil).Add(filters...)
	}
}
func (fil CompoundFilter) AndEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(EqualTo(field, value))
}
func (fil CompoundFilter) AndNotEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(NotEqualTo(field, value))
}
func (fil CompoundFilter) AndGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.And(GreaterThan(field, value))
}
func (fil CompoundFilter) AndNotGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.And(NotGreaterThan(field, value))
}
func (fil CompoundFilter) AndGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(GreaterThanOrEqualTo(field, value))
}
func (fil CompoundFilter) AndNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.And(NotGreaterThanOrEqualTo(field, value))
}

func (fil CompoundFilter) Or(filters ...Filter) CompoundFilter {
	if fil.Op == OpOr {
		return fil.Add(filters...)
	} else {
		return Or(fil).Add(filters...)
	}
}
func (fil CompoundFilter) OrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(EqualTo(field, value))
}
func (fil CompoundFilter) OrNotEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotEqualTo(field, value))
}
func (fil CompoundFilter) OrGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.Or(GreaterThan(field, value))
}
func (fil CompoundFilter) OrNotGreaterThan(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotGreaterThan(field, value))
}
func (fil CompoundFilter) OrGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(GreaterThanOrEqualTo(field, value))
}
func (fil CompoundFilter) OrNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter {
	return fil.Or(NotGreaterThanOrEqualTo(field, value))
}
