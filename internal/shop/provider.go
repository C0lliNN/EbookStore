package shop

import (
	catalog "github.com/c0llinn/ebook-store/internal/catalog/usecase"
	"github.com/c0llinn/ebook-store/internal/shop/client"
	"github.com/c0llinn/ebook-store/internal/shop/repository"
	"github.com/c0llinn/ebook-store/internal/shop/usecase"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	repository.NewOrderRepository,
	wire.Bind(new(usecase.Repository), new(repository.OrderRepository)),
	client.NewStripeClient,
	wire.Bind(new(usecase.PaymentClient), new(client.StripeClient)),
	wire.Bind(new(usecase.CatalogService), new(catalog.CatalogUseCase)),
	usecase.NewShopUseCase,
)
