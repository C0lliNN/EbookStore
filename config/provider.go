package config

import (
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	NewConnection,
	NewSNSService,
	NewBucket,
	NewS3Service,
)
