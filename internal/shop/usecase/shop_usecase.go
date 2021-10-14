package usecase

import (
	"github.com/c0llinn/ebook-store/internal/shop/model"
)
import catalog "github.com/c0llinn/ebook-store/internal/catalog/model"

type Repository interface {
	FindByQuery(query model.OrderQuery) (model.PaginatedOrders, error)
	FindByID(id string) (model.Order, error)
	Create(order *model.Order) error
	Update(order *model.Order) error
}

type PaymentClient interface {
	CreatePaymentIntentForOrder(order *model.Order) error
}

type CatalogService interface {
	FindBookByID(bookId string) (catalog.Book, error)
}

type ShopUseCase struct {
	repo Repository
	paymentClient PaymentClient
	catalogService CatalogService
}

func NewShopUseCase(repo Repository, paymentClient PaymentClient, catalogService CatalogService) ShopUseCase {
	return ShopUseCase{repo: repo, paymentClient: paymentClient, catalogService: catalogService}
}

func (u ShopUseCase) FindOrders(query model.OrderQuery) (model.PaginatedOrders, error) {
	return u.repo.FindByQuery(query)
}

func (u ShopUseCase) FindOrderByID(id string) (model.Order, error) {
	return u.repo.FindByID(id)
}

func (u ShopUseCase) CreateOrder(order *model.Order) error {
	book, err := u.catalogService.FindBookByID(order.BookID)
	if err != nil {
		return err
	}
	order.Total = int64(book.Price)

	if err = u.paymentClient.CreatePaymentIntentForOrder(order); err != nil {
		return err
	}

	return u.repo.Create(order)
}

func (u ShopUseCase) UpdateOrder(order *model.Order) error {
	return u.repo.Update(order)
}

func (u ShopUseCase) CompleteOrder(orderID string) error {
	order, err := u.repo.FindByID(orderID)
	if err != nil {
		return err
	}

	order.Complete()
	return u.repo.Update(&order)
}
