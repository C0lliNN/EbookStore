// Package query provides a way to create queries independent of the database technology.

package query

// ComparisonOperator is a string that represents a operators like =, !=, >, <, etc.
type ComparisonOperator string

const (
	Equal    ComparisonOperator = "="
	Match ComparisonOperator = "MATCH"
	NotEqual ComparisonOperator = "!="
)

// LogicalOperator is a string that represents a logical operator like AND, OR, etc.
type LogicalOperator string

const (
	and LogicalOperator = "AND"
	or  LogicalOperator = "OR"
)

// Condition is a struct that represents a condition in a query. Example: "name = 'John'"
type Condition struct {
	Field    string
	Operator ComparisonOperator
	Value    interface{}
}

type node struct {
	Value Condition
	Operator LogicalOperator
	Next  *node
}

// Query is a struct that represents a query. It is a linked list of nodes.
type Query struct {
	root *node
}

func New() *Query {
	return &Query{}
}

// And appends a condition to the query with a logical operator AND.
func (q *Query) And(condition Condition) *Query {
	q.appendCondition(condition, and)

	return q
}

// Or appends a condition to the query with a logical operator OR.
func (q *Query) Or(condition Condition) *Query {
	q.appendCondition(condition, or)

	return q
}

func (q *Query) appendCondition(condition Condition, logicalOperator LogicalOperator) {
	if q.Empty() {
		q.root = &node{
			Value: condition,
		}
		return
	}

	current := q.root

	for current.Next != nil {
		current = current.Next
	}

	current.Next =  &node{
		Operator: logicalOperator,
		Value: condition,
	}
}

func (q *Query) Empty() bool {
	return q.root == nil
}

type Iterator struct {
	current *node
}

func (q *Query) Iterator() *Iterator {
	return &Iterator{
		current: q.root,
	}
}

func (i *Iterator) HasNext() bool {
	return i.current != nil
}

func (i *Iterator) Next() (LogicalOperator, Condition) {
	if i.current == nil {
		return "", Condition{}
	}

	condition := i.current.Value
	operator := i.current.Operator
	i.current = i.current.Next

	return operator, condition
}
