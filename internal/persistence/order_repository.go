package persistence

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/c0llinn/ebook-store/internal/shop"
	"gorm.io/gorm"
	"strings"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepository{db: db}
}

func (r OrderRepository) FindByQuery(ctx context.Context, query shop.OrderQuery) (paginated shop.PaginatedOrders, err error) {
	conditions := r.createConditionsFromCriteria(query.CreateCriteria())
	result := r.db.Limit(query.Limit).Offset(query.Offset).Where(conditions).Order("created_at DESC").Find(&paginated.Orders)
	if err = result.Error; err != nil {
		log.Default().Error("error trying to find orders by query: ", err)
		return
	}

	var count int64
	if err = r.db.Model(&shop.Order{}).Where(conditions).Count(&count).Error; err != nil {
		log.Default().Error("error trying to count orders: ", err)
		return
	}

	paginated.Limit = query.Limit
	paginated.Offset = query.Offset
	paginated.TotalOrders = count

	return
}

func (r OrderRepository) createConditionsFromCriteria(criteria []shop.Criteria) string {
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

func (r OrderRepository) FindByID(ctx context.Context, id string) (order shop.Order, err error) {
	result := r.db.First(&order, "id = ?", id)
	if err = result.Error; err != nil {
		log.Default().Errorf("error when trying to find order with id %s: %v", id, err)
		err = &common.ErrEntityNotFound{Entity: "Order", Err: err}
	}

	return
}

func (r OrderRepository) Create(ctx context.Context, order *shop.Order) error {
	result := r.db.Create(order)
	if err := result.Error; err != nil {
		log.Default().Error("error trying to create an order: ", err)
		return err
	}

	return nil
}

func (r OrderRepository) Update(ctx context.Context, order *shop.Order) error {
	result := r.db.Updates(order).Where("id = ?", order.ID)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to update the order %s: %v", order.ID, err)
		return err
	}

	return nil
}
