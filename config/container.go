package config

import (
	"github.com/c0llinn/ebook-store/internal"
	"github.com/google/wire"
)

var Container = wire.NewSet(
	Provider,
	internal.NewRepository,
)
