package shop

import (
	"github.com/c0llinn/ebook-store/internal/shop/repository"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	repository.NewOrderRepository,
)
