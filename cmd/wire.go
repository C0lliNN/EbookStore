//+build wireinject

package main

import (
	"github.com/c0llinn/ebook-store/config/di"
	"github.com/c0llinn/ebook-store/internal/auth/repository"
	"github.com/google/wire"
)

func SetupApplication() repository.UserRepository {
	wire.Build(di.Container)
	return repository.UserRepository{}
}
