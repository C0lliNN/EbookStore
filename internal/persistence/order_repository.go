package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/c0llinn/ebook-store/internal/shop"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindByQuery(ctx context.Context, query shop.OrderQuery) (shop.PaginatedOrders, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	db := r.db.WithContext(ctx)
	conditions := r.createConditionsFromCriteria(query.CreateCriteria())

	paginated := shop.PaginatedOrders{}
	result := db.Limit(query.Limit).Offset(query.Offset).Where(conditions).Order("created_at DESC").Find(&paginated.Orders)
	if err := result.Error; err != nil {
		return shop.PaginatedOrders{}, fmt.Errorf("(FindByQuery) failed running select query: %w", err)
	}

	var count int64
	if err := db.Model(&shop.Order{}).Where(conditions).Count(&count).Error; err != nil {
		return shop.PaginatedOrders{}, fmt.Errorf("(FindByQuery) failed running count query: %w", err)
	}

	paginated.Limit = query.Limit
	paginated.Offset = query.Offset
	paginated.TotalOrders = count

	return paginated, nil
}

func (r *OrderRepository) createConditionsFromCriteria(criteria []shop.Criteria) string {
	conditions := make([]string, 0, len(criteria))

	for _, c := range criteria {
		if !c.IsEmpty() {
			if parsed, ok := c.Value.(string); ok {
				conditions = append(conditions, fmt.Sprintf("%s %s '%s'", c.Field, c.Operator, parsed))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s %s %v", c.Field, c.Operator, c.Value))
			}
		}
	}

	return strings.Join(conditions, " AND ")
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
