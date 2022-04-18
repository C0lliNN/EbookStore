package config

import (
	"github.com/c0llinn/ebook-store/internal/api"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/google/wire"
)

var Container = wire.NewSet(
	Provider,
	auth.Provider,
	catalog.Provider,
	shop.Provider,
	api.Provider,
)
