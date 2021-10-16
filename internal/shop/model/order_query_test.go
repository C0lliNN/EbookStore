//go:build unit
// +build unit

package model

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrderQueryTestSuite struct {
	suite.Suite
}

func TestOrderQueryRun(t *testing.T) {
	suite.Run(t, new(OrderQueryTestSuite))
}

func (s *OrderQueryTestSuite) TestCreateCriteria_EmptyStruct() {
	query := OrderQuery{}

	criteria := query.CreateCriteria()

	for _, c := range criteria {
		assert.True(s.T(), c.IsEmpty())
	}
}

func (s *OrderQueryTestSuite) TestCreateCriteria_WithStatus() {
	query := OrderQuery{Status: Paid}

	criteria := query.CreateCriteria()

	assert.Contains(s.T(), criteria, NewEqualCriteria("status", string(Paid)))
}

func (s *OrderQueryTestSuite) TestCreateCriteria_WithUserID() {
	query := OrderQuery{UserID: uuid.NewString()}

	criteria := query.CreateCriteria()

	assert.Contains(s.T(), criteria, NewEqualCriteria("user_id", query.UserID))
}

func (s *OrderQueryTestSuite) TestCreateCriteria_WithBookID() {
	query := OrderQuery{BookID: uuid.NewString()}

	criteria := query.CreateCriteria()

	assert.Contains(s.T(), criteria, NewEqualCriteria("book_id", query.BookID))
}

func (s *OrderQueryTestSuite) TestCreateCriteria_WithMultipleFields() {
	query := OrderQuery{BookID: uuid.NewString(), Status: Paid}

	criteria := query.CreateCriteria()

	assert.Contains(s.T(), criteria, NewEqualCriteria("book_id", query.BookID))
	assert.Contains(s.T(), criteria, NewEqualCriteria("status", string(Paid)))
}
