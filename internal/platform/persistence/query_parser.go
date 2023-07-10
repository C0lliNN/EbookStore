package persistence

import (
	"fmt"
	"strings"

	"github.com/ebookstore/internal/core/query"
)

var operatorMapping = map[query.ComparisonOperator]string{
	query.Equal: "=",
	query.Match: "ILIKE",
	query.NotEqual: "!=",
}

// parseQuery function responsible for parsing a query into a SQL string
func parseQuery(query query.Query) (string, []interface{}) {
	if query.Empty() {
		return "", nil
	}

	iterator := query.Iterator()
	var result strings.Builder
	var values []interface{}
	for iterator.HasNext() {
		operator, condition := iterator.Next()
		if len(operator) > 0 {
			result.WriteString(fmt.Sprintf(" %s ", operator))
		}

		field := condition.Field
		op := parseCondition(condition)
		value := parseValue(condition)

		result.WriteString(fmt.Sprintf("%s %s ?", field, op))
		values = append(values, value)
	}

	return result.String(), values
}

func parseCondition(condition query.Condition) string {
	if condition.Value == nil {
		if condition.Operator == query.Equal {
			return "IS"
		} else if condition.Operator == query.NotEqual {
			return "IS NOT"
		}
		return ""
	}

	return operatorMapping[condition.Operator]
}

func parseValue(condition query.Condition) interface{} {
	switch {
	case condition.Value == nil:
		return nil
	case condition.Operator == query.Match:
		return fmt.Sprintf("%%%s%%", condition.Value)
	default:
		return condition.Value
	}
}