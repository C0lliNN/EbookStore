package shop

import (
	"github.com/c0llinn/ebook-store/internal/shop/client"
	"github.com/c0llinn/ebook-store/internal/shop/delivery/http"
	"github.com/c0llinn/ebook-store/internal/shop/helper"
	"github.com/c0llinn/ebook-store/internal/shop/repository"
	"github.com/c0llinn/ebook-store/internal/shop/usecase"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	repository.NewOrderRepository,
	wire.Bind(new(usecase.Repository), new(repository.OrderRepository)),
	client.NewStripeClient,
	wire.Bind(new(usecase.PaymentClient), new(client.StripeClient)),
	usecase.NewShopUseCase,
	wire.Bind(new(http.ShopService), new(usecase.ShopUseCase)),
	helper.NewIDGenerator,
	wire.Bind(new(http.IDGenerator), new(helper.IDGenerator)),
	http.NewShopHandler,
)
