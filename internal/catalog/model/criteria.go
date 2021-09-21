package model

import (
	"reflect"
)

type Criteria struct {
	Field string
	Operator string
	Value interface{}
}

func (c Criteria) IsEmpty() bool {
	return c.Value == nil || reflect.ValueOf(c.Value).IsZero()
}

func NewEqualCriteria(field string, value interface{}) Criteria {
	return Criteria{field, "=", value}
}

func NewILikeCriteria(field, value string) Criteria {
	if value == "" {
		return Criteria{}
	}

	return Criteria{field, "ILIKE", "%" + value + "%"}
}