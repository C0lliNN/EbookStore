package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/test"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type CartRepositoryTestSuite struct {
	suite.Suite
	repo      *CartRepository
	client    *redisclient.Client
	container *test.RedisContainer
}

func (s *CartRepositoryTestSuite) SetupSuite() {
	ctx := context.TODO()

	var err error
	s.container, err = test.NewRedisContainer(ctx)
	s.Require().NoError(err)

	s.client = redisclient.NewClient(&redisclient.Options{
		Addr: s.container.Endpoint,
	})
	s.repo = NewCartRepository(s.client, time.Second*10)
}

func (s *CartRepositoryTestSuite) TearDownSuite() {
	ctx := context.TODO()

	s.Require().NoError(s.container.Terminate(ctx))
}

func (s *CartRepositoryTestSuite) TearDownTest() {
	ctx := context.TODO()

	s.Require().NoError(s.client.FlushDB(ctx).Err())
}

func (s *CartRepositoryTestSuite) TestFindByUserID() {
	ctx := context.TODO()

	cart := &shop.Cart{
		UserID: "1",
		Items: []shop.Item{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		},
	}

	s.Require().NoError(s.repo.Save(ctx, cart))

	result, err := s.repo.FindByUserID(ctx, "1")
	s.NoError(err)
	s.Equal(cart, result)
}

func (s *CartRepositoryTestSuite) TestFindByUserID_NotFound() {
	ctx := context.TODO()

	result, err := s.repo.FindByUserID(ctx, "1")
	s.Equal(&ErrEntityNotFound{entity: "cart"}, err)
	s.Nil(result)
}

func (s *CartRepositoryTestSuite) TestSave() {
	ctx := context.TODO()

	cart := &shop.Cart{
		UserID: "1",
		Items: []shop.Item{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		},
	}

	s.NoError(s.repo.Save(ctx, cart))

	result, err := s.repo.FindByUserID(ctx, "1")
	s.NoError(err)
	s.Equal(cart, result)
}

func (s *CartRepositoryTestSuite) TestDeleteByUserID() {
	ctx := context.TODO()

	cart := &shop.Cart{
		UserID: "1",
		Items: []shop.Item{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		},
	}

	s.Require().NoError(s.repo.Save(ctx, cart))
	s.Require().NoError(s.repo.DeleteByUserID(ctx, "1"))

	result, err := s.repo.FindByUserID(ctx, "1")
	s.Error(err)
	s.Nil(result)
}

func TestCartRepository(t *testing.T) {
	suite.Run(t, new(CartRepositoryTestSuite))
}
