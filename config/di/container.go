package di

import (
	"github.com/c0llinn/ebook-store/config"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/google/wire"
)

var Container = wire.NewSet(
	config.Provider,
	auth.Provider,
)
