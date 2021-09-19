package catalog

import (
	"github.com/c0llinn/ebook-store/internal/catalog/delivery/http"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	http.NewCatalogHandler,
)
