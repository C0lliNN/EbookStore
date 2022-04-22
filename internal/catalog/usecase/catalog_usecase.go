package usecase

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"io"
)

type Repository interface {
	FindByQuery(ctx context.Context, query model.BookQuery) (paginated model.PaginatedBooks, err error)
	FindByID(ctx context.Context, id string) (book model.Book, err error)
	Create(ctx context.Context, book *model.Book) error
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id string) error
}

type StorageClient interface {
	GeneratePreSignedUrl(ctx context.Context, key string) (string, error)
	SaveFile(ctx context.Context, key string, contentType string, content io.ReadSeeker) error
	RetrieveFile(ctx context.Context, key string) (io.ReadCloser, error)
}

type FilenameGenerator interface {
	NewUniqueName(filename string) string
}

type CatalogUseCase struct {
	repo              Repository
	storageClient     StorageClient
	filenameGenerator FilenameGenerator
}

func NewCatalogUseCase(repo Repository, storageClient StorageClient, filenameGenerator FilenameGenerator) CatalogUseCase {
	return CatalogUseCase{repo: repo, storageClient: storageClient, filenameGenerator: filenameGenerator}
}

func (u CatalogUseCase) FindBooks(ctx context.Context, query model.BookQuery) (paginatedBooks model.PaginatedBooks, err error) {
	paginatedBooks, err = u.repo.FindByQuery(ctx, query)
	if err != nil {
		return
	}

	for i := range paginatedBooks.Books {
		imageKey := paginatedBooks.Books[i].PosterImageBucketKey

		var url string
		url, err = u.storageClient.GeneratePreSignedUrl(ctx, imageKey)
		if err != nil {
			return
		}

		paginatedBooks.Books[i].SetPosterImageLink(url)
	}

	return
}

func (u CatalogUseCase) FindBookByID(ctx context.Context, id string) (book model.Book, err error) {
	book, err = u.repo.FindByID(ctx, id)
	if err != nil {
		return
	}

	imageKey := book.PosterImageBucketKey
	url, err := u.storageClient.GeneratePreSignedUrl(ctx, imageKey)
	if err != nil {
		return
	}

	book.SetPosterImageLink(url)
	return
}

func (u CatalogUseCase) GetBookContent(ctx context.Context, id string) (io.ReadCloser, error) {
	book, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.storageClient.RetrieveFile(ctx, book.ContentBucketKey)
}

func (u CatalogUseCase) CreateBook(ctx context.Context, book *model.Book, posterImage io.ReadSeeker, bookContent io.ReadSeeker) error {
	posterImageKey := u.filenameGenerator.NewUniqueName("poster_" + book.Title)
	contentKey := u.filenameGenerator.NewUniqueName("content_" + book.Title)

	if err := u.storageClient.SaveFile(ctx, posterImageKey, "image/jpeg", posterImage); err != nil {
		return err
	}

	if err := u.storageClient.SaveFile(ctx, contentKey, "application/pdf", bookContent); err != nil {
		return err
	}

	book.PosterImageBucketKey = posterImageKey
	book.ContentBucketKey = contentKey

	url, err := u.storageClient.GeneratePreSignedUrl(ctx, posterImageKey)
	if err != nil {
		return err
	}
	book.SetPosterImageLink(url)

	return u.repo.Create(ctx, book)
}

func (u CatalogUseCase) UpdateBook(ctx context.Context, book *model.Book) error {
	return u.repo.Update(ctx, book)
}

func (u CatalogUseCase) DeleteBook(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
