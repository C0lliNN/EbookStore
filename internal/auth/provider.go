package auth

import (
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	NewUserRepository,
)
