//+build wireinject

package main

import (
	"github.com/c0llinn/ebook-store/config/di"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/google/wire"
)

func SetupApplication() auth.UserRepository {
	wire.Build(di.Container)
	return auth.UserRepository{}
}
