package config

import (
	"github.com/c0llinn/ebook-store/config/aws"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	db.NewConnection,
	aws.NewSNSService,
)
