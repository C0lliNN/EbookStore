package config

import (
	"github.com/c0llinn/ebook-store/internal"
	"github.com/google/wire"
)

var Container = wire.NewSet(internal.NewBar, internal.NewFoo, Provider)
