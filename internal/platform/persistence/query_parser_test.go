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
		expected string
	} {
		{
			name: "when query is empty, then it should return an empty string",
			query: *query.New(),
			expected: "",
		},
		{
			name: "when query has only one condition, then it should return a string with the condition",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.Equal, Value: "value"}),
			expected: "title = 'value'",
		},
		{
			name: "when query has two conditions, then it should return a string with the conditions",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.Equal, Value: "value"}).
				And(query.Condition{Field: "author", Operator: query.Equal, Value: "author"}),
			expected: "title = 'value' AND author = 'author'",
		},
		{
			name: "when query has three conditions, then it should return a string with the conditions",
			query: *query.New().And(query.Condition{Field: "title", Operator: query.NotEqual, Value: nil}).
				And(query.Condition{Field: "author", Operator: query.Equal, Value: "author"}).
				And(query.Condition{Field: "price", Operator: query.Equal, Value: 10}),
			expected: "title IS NOT NULL AND author = 'author' AND price = 10",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseQuery(tc.query)
			assert.Equal(t, actual, tc.expected)
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
		expected string
	} {
		{
			name: "when operator is match, then it should return a string in the format ILIKE",
			condition: query.Condition{Field: "title", Operator: query.Match, Value: "value"},
			expected: "'%value%'",
		},
		{
			name: "when value is a string, then it should return a string in the format 'value'",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: "value"},
			expected: "'value'",
		},
		{
			name: "when value is nil, then it should return NULL",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: nil},
			expected: "NULL",
		},
		{
			name: "when value does not fit into the other cases, then it should return a string representation of the value",
			condition: query.Condition{Field: "title", Operator: query.Equal, Value: 1},
			expected: "1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseValue(tc.condition)
			assert.Equal(t, actual, tc.expected)
		})
	}
}