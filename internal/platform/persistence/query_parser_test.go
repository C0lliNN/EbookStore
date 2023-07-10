package persistence

import (
	"testing"

	"github.com/ebookstore/internal/core/query"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    query.Query
		expectedQuery string
		expectedValues []interface{}
	} {
		{
			name: "when query is empty, then it should return an empty string",
			query: *query.New(),
			expectedQuery: "",
			expectedValues: nil,
		},
		{
			name: "when query has ILIKE condition, then it should return a string with ILIKE and the value should have %",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.Match, Value: "value"}),
			expectedQuery: "title ILIKE ?",
			expectedValues: []interface{}{"%value%"},
		},
		{
			name: "when query has only one condition, then it should return a string with the condition",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.Equal, Value: "value"}),
			expectedQuery: "title = ?",
			expectedValues: []interface{}{"value"},
		},
		{
			name: "when query has two conditions, then it should return a string with the conditions",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.Equal, Value: "value"}).
				And(query.Condition{Field: "author", Operator: query.Equal, Value: "author"}),
			expectedQuery: "title = ? AND author = ?",
			expectedValues: []interface{}{"value", "author"},
		},
		{
			name: "when query has three conditions, then it should return a string with the conditions",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.NotEqual, Value: nil}).
				And(query.Condition{Field: "author", Operator: query.Equal, Value: "author"}).
				And(query.Condition{Field: "price", Operator: query.Equal, Value: 10}),
			expectedQuery: "title IS NOT ? AND author = ? AND price = ?",
			expectedValues: []interface{}{nil, "author", 10},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualQuery, actualValues := parseQuery(tc.query)
			assert.Equal(t, actualValues, tc.expectedValues)
			assert.Equal(t, actualQuery, tc.expectedQuery)
		})
	}
}

func TestParseCondition(t *testing.T) {
	tests := []struct {
		name     string
		condition query.Condition
		expected string
	} {
		{
			name: "when operator is equal and value is nil, then it should return IS",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: nil},
			expected: "IS",
		},
		{
			name: "when operator is not equal and value is nil, then it should return IS NOT",
			condition: query.Condition{Field: "title", Operator: query.NotEqual, Value: nil},
			expected: "IS NOT",
		},
		{
			name: "when operator is equal and value is not nil, then it should return =",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: "value"},
			expected: "=",
		},
		{
			name: "when operator is not equal and value is not nil, then it should return !=",
			condition: query.Condition{Field: "title", Operator: query.NotEqual, Value: "value"},
			expected: "!=",
		},
		{
			name: "when operator is match, then it should return ILIKE",
			condition: query.Condition{Field: "title", Operator: query.Match, Value: "value"},
			expected: "ILIKE",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseCondition(tc.condition)
			assert.Equal(t, actual, tc.expected)
		})
	}
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		condition query.Condition
		expected interface{}
	} {
		{
			name: "when value is a string, then it should return the string",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: "value"},
			expected: "value",
		},
		{
			name: "when the operator is match, then it should return the string surrounded with %",
			condition: query.Condition{Field: "title", Operator: query.Match, Value: "value"},
			expected: "%value%",
		},
		{
			name: "when value is an int, then it should return the int",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: 10},
			expected: 10,
		},
		{
			name: "when value is nil, then it should return nil",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: nil},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseValue(tc.condition)
			assert.Equal(t, actual, tc.expected)
		})
	}
}