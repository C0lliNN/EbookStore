package config

import (
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/env"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/google/wire"
)

func init() {
	env.InitConfiguration()
	log.InitLogger()
}

var Provider = wire.NewSet(
	db.NewConnection,
)
