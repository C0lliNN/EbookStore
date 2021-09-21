// +build unit

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBookQuery_CreateCriteria_WithEmptyData(t *testing.T) {
	query := BookQuery{}

	expected := []Criteria{
		{},
		{},
		{Field: "author_name", Operator: "=", Value: ""},
	}

	actual := query.CreateCriteria()

	assert.Equal(t, expected, actual)
}

func TestBookQuery_CreateCriteria_WithTitle(t *testing.T) {
	query := BookQuery{Title: "some title"}

	expected := []Criteria{
		{Field: "title", Operator: "ILIKE", Value: "%some title%"},
		{},
		{Field: "author_name", Operator: "=", Value: ""},
	}

	actual := query.CreateCriteria()

	assert.Equal(t, expected, actual)
}

func TestBookQuery_CreateCriteria_WithDescription(t *testing.T) {
	query := BookQuery{Description: "some description"}

	expected := []Criteria{
		{},
		{Field: "description", Operator: "ILIKE", Value: "%some description%"},
		{Field: "author_name", Operator: "=", Value: ""},
	}

	actual := query.CreateCriteria()

	assert.Equal(t, expected, actual)
}

func TestBookQuery_CreateCriteria_WithAuthorName(t *testing.T) {
	query := BookQuery{AuthorName: "some name"}

	expected := []Criteria{
		{},
		{},
		{Field: "author_name", Operator: "=", Value: "some name"},
	}

	actual := query.CreateCriteria()

	assert.Equal(t, expected, actual)
}

func TestBookQuery_CreateCriteria_WithAllFields(t *testing.T) {
	query := BookQuery{Title: "some title", Description: "some description", AuthorName: "some name"}

	expected := []Criteria{
		{Field: "title", Operator: "ILIKE", Value: "%some title%"},
		{Field: "description", Operator: "ILIKE", Value: "%some description%"},
		{Field: "author_name", Operator: "=", Value: "some name"},
	}

	actual := query.CreateCriteria()

	assert.Equal(t, expected, actual)
}
