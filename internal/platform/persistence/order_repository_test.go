package persistence_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/platform/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *persistence.OrderRepository
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	s.RepositoryTestSuite.SetupSuite()

	s.repo = persistence.NewOrderRepository(s.db)
}

func TestOrderRepositoryRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TearDownTest() {
	s.db.Delete(&shop.Order{}, "1 = 1")
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithEmptyQuery() {
	order1 := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order2 := shop.Order{
		ID:     "some-id2",
		Status: shop.Pending,
		Total:  4000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order3 := shop.Order{
		ID:     "some-id3",
		Status: shop.Pending,
		Total:  3500,
		BookID: "book-id",
		UserID: "user-id",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(ctx, shop.OrderQuery{})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(3), actual.TotalOrders)
	assert.Equal(s.T(), order3.ID, actual.Orders[0].ID)
	assert.Equal(s.T(), order2.ID, actual.Orders[1].ID)
	assert.Equal(s.T(), order1.ID, actual.Orders[2].ID)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithLimit() {
	order1 := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order2 := shop.Order{
		ID:     "some-id2",
		Status: shop.Pending,
		Total:  4000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order3 := shop.Order{
		ID:     "some-id3",
		Status: shop.Pending,
		Total:  3500,
		BookID: "book-id",
		UserID: "user-id",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(ctx, shop.OrderQuery{Limit: 1})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 1, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(3), actual.TotalOrders)
	assert.Equal(s.T(), order3.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithOffset() {
	order1 := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order2 := shop.Order{
		ID:     "some-id2",
		Status: shop.Pending,
		Total:  4000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order3 := shop.Order{
		ID:     "some-id3",
		Status: shop.Pending,
		Total:  3500,
		BookID: "book-id",
		UserID: "user-id",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(ctx, shop.OrderQuery{Offset: 1})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 1, actual.Offset)
	assert.Equal(s.T(), int64(3), actual.TotalOrders)
	assert.Equal(s.T(), order2.ID, actual.Orders[0].ID)
	assert.Equal(s.T(), order1.ID, actual.Orders[1].ID)
	assert.Len(s.T(), actual.Orders, 2)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithStatus() {
	order1 := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order2 := shop.Order{
		ID:     "some-id2",
		Status: shop.Pending,
		Total:  4000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order3 := shop.Order{
		ID:     "some-id3",
		Status: shop.Pending,
		Total:  3500,
		BookID: "book-id",
		UserID: "user-id",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(ctx, shop.OrderQuery{Status: shop.Paid})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(1), actual.TotalOrders)
	assert.Equal(s.T(), order1.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithBookID() {
	order1 := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id2",
		UserID: "user-id",
	}
	order2 := shop.Order{
		ID:     "some-id2",
		Status: shop.Pending,
		Total:  4000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order3 := shop.Order{
		ID:     "some-id3",
		Status: shop.Pending,
		Total:  3500,
		BookID: "book-id",
		UserID: "user-id",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(ctx, shop.OrderQuery{BookID: order1.BookID})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(1), actual.TotalOrders)
	assert.Equal(s.T(), order1.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithUserID() {
	order1 := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id2",
		UserID: "user-id2",
	}
	order2 := shop.Order{
		ID:     "some-id2",
		Status: shop.Pending,
		Total:  4000,
		BookID: "book-id",
		UserID: "user-id",
	}
	order3 := shop.Order{
		ID:     "some-id3",
		Status: shop.Pending,
		Total:  3500,
		BookID: "book-id",
		UserID: "user-id",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(ctx, shop.OrderQuery{UserID: order1.UserID})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(1), actual.TotalOrders)
	assert.Equal(s.T(), order1.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByID_Successfully() {
	order := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id2",
		UserID: "user-id2",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByID(ctx, order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), order.ID, actual.ID)
}

func (s *OrderRepositoryTestSuite) TestFindByID_WithError() {
	ctx := context.TODO()

	_, err := s.repo.FindByID(ctx, uuid.NewString())
	assert.IsType(s.T(), &persistence.ErrEntityNotFound{}, errors.Unwrap(err))
}

func (s *OrderRepositoryTestSuite) TestCreate_Successfully() {
	order := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id2",
		UserID: "user-id2",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order)
	assert.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(ctx, order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), order.ID, persisted.ID)
}

func (s *OrderRepositoryTestSuite) TestCreate_WithError() {
	order := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id2",
		UserID: "user-id2",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order)
	assert.Nil(s.T(), err)

	err = s.repo.Create(ctx, &order)
	assert.NotNil(s.T(), err)
}

func (s *OrderRepositoryTestSuite) TestUpdate_Successfully() {
	order := shop.Order{
		ID:     "some-id1",
		Status: shop.Paid,
		Total:  5000,
		BookID: "book-id2",
		UserID: "user-id2",
	}

	ctx := context.TODO()

	err := s.repo.Create(ctx, &order)
	assert.Nil(s.T(), err)

	order.Status = shop.Paid
	err = s.repo.Update(ctx, &order)
	assert.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(ctx, order.ID)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), shop.Paid, persisted.Status)
}
