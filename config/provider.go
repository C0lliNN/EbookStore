package config

import "github.com/google/wire"

func init() {
	InitConfiguration()
	InitLogger()
}

var Provider = wire.NewSet(
	NewConnection,
)
