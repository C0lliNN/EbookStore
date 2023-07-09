package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	expected := &Query{}
	actual := New()

	assert.Equal(t, expected, actual)
}

func TestQuery_And(t *testing.T) {
	tests := []struct {
		name     string
		query    *Query
		field    string
		operator ComparisonOperator
		value    interface{}
		expected *Query
	}{
		{
			name:     "when query is empty, then it should add a new node",
			query:    New(),
			field:    "book_id",
			operator: Equal,
			value:    "id",
			expected: &Query{
				root: &node{
					Value: Condition{
						Field:    "book_id",
						Operator: Equal,
						Value:    "id",
					},
				},
			},
		},
		{
			name:     "when query has only the root, then it should add a new node",
			query:    New().And(Condition{"book_id", Equal, "id"}),
			field:    "title",
			operator: Equal,
			value:    "value",
			expected: &Query{
				root: &node{
					Value: Condition{
						Field:    "book_id",
						Operator: Equal,
						Value:    "id",
					},
					Next:  &node{
							Operator: and,
							Value: Condition{
								Field:    "title",
								Operator: Equal,
								Value:    "value",
							},
						},
					
				},
			},
		},
		{
			name:     "when query has more than one node, then it should add a new node",
			query:    New().And(Condition{"book_id", Equal, "id"}).And(Condition{"title", Equal, "value"}),
			field:    "author",
			operator: Equal,
			value:    "value",
			expected: &Query{
				root: &node{
					Value: Condition{
						Field:    "book_id",
						Operator: Equal,
						Value:    "id",
					},
					Next:  &node{
							Operator: and,
							Value: Condition{
								Field:    "title",
								Operator: Equal,
								Value:    "value",
							},
							Next: &node{
								Operator: and,
									Value: Condition{
										Field:    "author",
										Operator: Equal,
										Value:    "value",
									},
								},
							},
						},
				
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.query.And(Condition{tc.field, tc.operator, tc.value})

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestQuery_Or(t *testing.T) {
	tests := []struct {
		name     string
		query    *Query
		field    string
		operator ComparisonOperator
		value    interface{}
		expected *Query
	}{
		{
			name:     "when query is empty, then it should add a new node",
			query:    New(),
			field:    "book_id",
			operator: Equal,
			value:    "id",
			expected: &Query{
				root: &node{
					Value: Condition{
						Field:    "book_id",
						Operator: Equal,
						Value:    "id",
					},
				},
			},
		},
		{
			name:     "when query has only the root, then it should add a new node",
			query:    New().Or(Condition{"book_id", Equal, "id"}),
			field:    "title",
			operator: Equal,
			value:    "value",
			expected: &Query{
				root: &node{
					Value: Condition{
						Field:    "book_id",
						Operator: Equal,
						Value:    "id",
					},
					Next:  &node{
							Operator: or,
							Value: Condition{
								Field:    "title",
								Operator: Equal,
								Value:    "value",
							},
						},
					},
				
			},
		},
		{
			name:     "when query has more than one node, then it should add a new node",
			query:    New().Or(Condition{"book_id", Equal, "id"}).Or(Condition{"title", Equal, "value"}),
			field:    "author",
			operator: Equal,
			value:    "value",
			expected: &Query{
				root: &node{
					Value: Condition{
						Field:    "book_id",
						Operator: Equal,
						Value:    "id",
					},
					Next: &node{
						Operator: or,
							Value: Condition{
								Field:    "title",
								Operator: Equal,
								Value:    "value",
							},
							Next: &node{
								Operator: or,
									Value: Condition{
										Field:    "author",
										Operator: Equal,
										Value:    "value",
									},
								},
							},
						},
					},
				
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.query.Or(Condition{tc.field, tc.operator, tc.value})

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestQuery_Empty(t *testing.T) {
	tests := []struct {
		name     string
		query    *Query
		expected bool
	}{
		{
			name:     "when query has only the root, then it should return false",
			query:    New().And(Condition{"book_id", Equal, "id"}),
			expected: false,
		},
		{
			name:     "when query has more than one node, then it should return false",
			query:    New().And(Condition{"book_id", Equal, "id"}).And(Condition{"title", Equal, "value"}),
			expected: false,
		},
		{
			name:     "when query has no nodes, then it should return true",
			query:    New(),
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.query.Empty()

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestQuery_Iterator(t *testing.T) {
	query := New().And(Condition{"book_id", Equal, "id"}).And(Condition{"title", Equal, "value"})

	expected := &Iterator{current: query.root}
	actual := query.Iterator()

	assert.Equal(t, expected, actual)
}

func TestIterator_HasNext(t *testing.T) {
	query := New()
	assert.False(t, query.Iterator().HasNext())

	assert.True(t, query.And(Condition{"book_id", Equal, "id"}).Iterator().HasNext())
}

func TestIterator_Next(t *testing.T) {
	query := New().And(Condition{"book_id", Equal, "id"}).Or(Condition{"title", Equal, "value"})
	iterator := query.Iterator()

	operator, condition := iterator.Next()
	assert.Empty(t, operator)
	assert.Equal(t, Condition{"book_id", Equal, "id"}, condition)

	operator, condition = iterator.Next()
	assert.Equal(t, or, operator)
	assert.Equal(t, Condition{"title", Equal, "value"}, condition)
}