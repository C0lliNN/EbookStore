package catalog

import (
	"github.com/c0llinn/ebook-store/internal/catalog/delivery/http"
	"github.com/c0llinn/ebook-store/internal/catalog/helper"
	"github.com/c0llinn/ebook-store/internal/catalog/repository"
	"github.com/c0llinn/ebook-store/internal/catalog/storage"
	"github.com/c0llinn/ebook-store/internal/catalog/usecase"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	repository.NewBookRepository,
	wire.Bind(new(usecase.Repository), new(repository.BookRepository)),
	storage.NewS3Client,
	wire.Bind(new(usecase.StorageClient), new(storage.S3Client)),
	helper.NewFilenameGenerator,
	wire.Bind(new(usecase.FilenameGenerator), new(helper.FilenameGenerator)),
	usecase.NewCatalogUseCase,
	wire.Bind(new(http.Service), new(usecase.CatalogUseCase)),
	http.NewCatalogHandler,
)
