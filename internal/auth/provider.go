package auth

import (
	"github.com/c0llinn/ebook-store/internal/auth/repository"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	repository.NewUserRepository,
)
