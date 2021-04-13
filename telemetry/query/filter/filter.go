package filter

import "github.com/nerdynick/confluent-cloud-metrics-go-sdk/telemetry/labels"

type Filter interface {
	And(filters ...Filter) CompoundFilter
	AndEqualTo(field labels.Label, value string) CompoundFilter
	AndNotEqualTo(field labels.Label, value string) CompoundFilter
	AndGreaterThan(field labels.Label, value string) CompoundFilter
	AndNotGreaterThan(field labels.Label, value string) CompoundFilter
	AndGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter
	AndNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter
	Or(filters ...Filter) CompoundFilter
	OrEqualTo(field labels.Label, value string) CompoundFilter
	OrNotEqualTo(field labels.Label, value string) CompoundFilter
	OrGreaterThan(field labels.Label, value string) CompoundFilter
	OrNotGreaterThan(field labels.Label, value string) CompoundFilter
	OrGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter
	OrNotGreaterThanOrEqualTo(field labels.Label, value string) CompoundFilter
}

func NotAnyOf(filters ...Filter) UnaryFilter {
	return Not(Or(filters...))
}
func AnyOf(filters ...Filter) CompoundFilter {
	return Or(filters...)
}
func OneOf(filters ...Filter) CompoundFilter {
	return Or(filters...)
}
func Or(filters ...Filter) CompoundFilter {
	return CompoundFilter{
		Op:      OpOr,
		Filters: filters,
	}
}

func AllOf(filters ...Filter) CompoundFilter {
	return And(filters...)
}
func And(filters ...Filter) CompoundFilter {
	return CompoundFilter{
		Op:      OpAnd,
		Filters: filters,
	}
}

func Not(filter Filter) UnaryFilter {
	return UnaryFilter{
		Op:        OpNot,
		SubFilter: filter,
	}
}

func NotEqualTo(field labels.Label, value string) UnaryFilter {
	return Not(EqualTo(field, value))
}
func EqualTo(field labels.Label, value string) FieldFilter {
	return FieldFilter{
		Op:    OpEq,
		Field: field,
		Value: value,
	}
}

func NotGreaterThan(field labels.Label, value string) UnaryFilter {
	return Not(GreaterThan(field, value))
}
func GreaterThan(field labels.Label, value string) FieldFilter {
	return FieldFilter{
		Op:    OpGt,
		Field: field,
		Value: value,
	}
}
func NotGreaterThanOrEqualTo(field labels.Label, value string) UnaryFilter {
	return Not(GreaterThanOrEqualTo(field, value))
}
func GreaterThanOrEqualTo(field labels.Label, value string) FieldFilter {
	return FieldFilter{
		Op:    OpGte,
		Field: field,
		Value: value,
	}
}
