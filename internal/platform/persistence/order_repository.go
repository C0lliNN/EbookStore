package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ebookstore/internal/core/query"
	"github.com/ebookstore/internal/core/shop"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindByQuery(ctx context.Context, q query.Query, p query.Page) (shop.PaginatedOrders, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	db := r.db.WithContext(ctx)
	conditions := parseQuery(q)

	paginated := shop.PaginatedOrders{}
	result := db.Limit(p.Size).Offset(p.Offset()).Where(conditions).Order("created_at DESC").Find(&paginated.Orders)
	if err := result.Error; err != nil {
		return shop.PaginatedOrders{}, fmt.Errorf("(FindByQuery) failed running select query: %w", err)
	}

	var count int64
	if err := db.Model(&shop.Order{}).Where(conditions).Count(&count).Error; err != nil {
		return shop.PaginatedOrders{}, fmt.Errorf("(FindByQuery) failed running count query: %w", err)
	}

	paginated.Limit = p.Size
	paginated.Offset = p.Offset()
	paginated.TotalOrders = count

	return paginated, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (shop.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	order := shop.Order{}
	result := r.db.WithContext(ctx).First(&order, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = &ErrEntityNotFound{entity: "order"}
		}

		return shop.Order{}, fmt.Errorf("(FindByID) failed running select query: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *shop.Order) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Create(order)
	if err := result.Error; err != nil {
		return fmt.Errorf("(Create) failed running insert statement: %w", err)
	}

	return nil
}

func (r *OrderRepository) Update(ctx context.Context, order *shop.Order) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Save(order).Where("id = ?", order.ID)
	if err := result.Error; err != nil {
		return fmt.Errorf("(Update) failed running update statement: %w", err)
	}

	return nil
}
