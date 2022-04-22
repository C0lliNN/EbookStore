package usecase

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"io"
)
import catalog "github.com/c0llinn/ebook-store/internal/catalog/model"

type Repository interface {
	FindByQuery(ctx context.Context, query model.OrderQuery) (model.PaginatedOrders, error)
	FindByID(ctx context.Context, id string) (model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
}

type PaymentClient interface {
	CreatePaymentIntentForOrder(ctx context.Context, order *model.Order) error
}

type CatalogService interface {
	FindBookByID(ctx context.Context, bookId string) (catalog.Book, error)
	GetBookContent(ctx context.Context, bookId string) (io.ReadCloser, error)
}

type ShopUseCase struct {
	repo           Repository
	paymentClient  PaymentClient
	catalogService CatalogService
}

func NewShopUseCase(repo Repository, paymentClient PaymentClient, catalogService CatalogService) ShopUseCase {
	return ShopUseCase{repo: repo, paymentClient: paymentClient, catalogService: catalogService}
}

func (u ShopUseCase) FindOrders(ctx context.Context, query model.OrderQuery) (model.PaginatedOrders, error) {
	return u.repo.FindByQuery(ctx, query)
}

func (u ShopUseCase) FindOrderByID(ctx context.Context, id string) (model.Order, error) {
	return u.repo.FindByID(ctx, id)
}

func (u ShopUseCase) CreateOrder(ctx context.Context, order *model.Order) error {
	book, err := u.catalogService.FindBookByID(ctx, order.BookID)
	if err != nil {
		return err
	}
	order.Total = int64(book.Price)

	if err = u.paymentClient.CreatePaymentIntentForOrder(ctx, order); err != nil {
		return err
	}

	return u.repo.Create(ctx, order)
}

func (u ShopUseCase) UpdateOrder(ctx context.Context, order *model.Order) error {
	return u.repo.Update(ctx, order)
}

func (u ShopUseCase) CompleteOrder(ctx context.Context, orderID string) error {
	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Complete()
	return u.repo.Update(ctx, &order)
}

func (u ShopUseCase) DownloadOrder(ctx context.Context, orderID string) (io.ReadCloser, error) {
	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != model.Paid {
		return nil, &common.ErrOrderNotPaid{Err: fmt.Errorf("only books from paid order can be downloaded")}
	}

	return u.catalogService.GetBookContent(ctx, order.BookID)
}
