//+build wireinject

package app

import (
	"github.com/c0llinn/ebook-store/config"
	"github.com/c0llinn/ebook-store/internal"
	"github.com/google/wire"
)

func SetupApplication() internal.Repository {
	wire.Build(config.Container)
	return internal.Repository{}
}
