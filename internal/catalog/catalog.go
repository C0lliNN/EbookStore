package catalog

import (
	"context"
	"io"
)

type Repository interface {
	FindByQuery(ctx context.Context, query BookQuery) (paginated PaginatedBooks, err error)
	FindByID(ctx context.Context, id string) (book Book, err error)
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
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

type IDGenerator interface {
	NewID() string
}

type Validator interface {
	Validate(i interface{}) error
}

type Config struct {
	Repository        Repository
	StorageClient     StorageClient
	FilenameGenerator FilenameGenerator
	IDGenerator       IDGenerator
	Validator         Validator
}

type Catalog struct {
	Config
}

func New(c Config) *Catalog {
	return &Catalog{Config: c}
}

func (c *Catalog) FindBooks(ctx context.Context, request SearchBooks) (PaginatedBooksResponse, error) {
	query := request.BookQuery()

	paginatedBooks, err := c.Repository.FindByQuery(ctx, query)
	if err != nil {
		return PaginatedBooksResponse{}, err
	}

	for i := range paginatedBooks.Books {
		imageKey := paginatedBooks.Books[i].PosterImageBucketKey

		var url string
		url, err = c.StorageClient.GeneratePreSignedUrl(ctx, imageKey)
		if err != nil {
			return PaginatedBooksResponse{}, err
		}

		paginatedBooks.Books[i].SetPosterImageLink(url)
	}

	return NewPaginatedBooksResponse(paginatedBooks), nil
}

func (c *Catalog) FindBookByID(ctx context.Context, id string) (BookResponse, error) {
	book, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		return BookResponse{}, err
	}

	imageKey := book.PosterImageBucketKey
	url, err := c.StorageClient.GeneratePreSignedUrl(ctx, imageKey)
	if err != nil {
		return BookResponse{}, err
	}

	book.SetPosterImageLink(url)
	return NewBookResponse(book), nil
}

func (c *Catalog) GetBookContent(ctx context.Context, id string) (io.ReadCloser, error) {
	book, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return c.StorageClient.RetrieveFile(ctx, book.ContentBucketKey)
}

func (c *Catalog) CreateBook(ctx context.Context, request CreateBook) (BookResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		return BookResponse{}, err
	}

	book := request.Book(c.IDGenerator.NewID())
	posterImageKey := c.FilenameGenerator.NewUniqueName("poster_" + book.Title)
	contentKey := c.FilenameGenerator.NewUniqueName("content_" + book.Title)

	if err := c.StorageClient.SaveFile(ctx, posterImageKey, "image/jpeg", request.PosterImage); err != nil {
		return BookResponse{}, err
	}

	if err := c.StorageClient.SaveFile(ctx, contentKey, "application/pdf", request.BookContent); err != nil {
		return BookResponse{}, err
	}

	book.PosterImageBucketKey = posterImageKey
	book.ContentBucketKey = contentKey

	url, err := c.StorageClient.GeneratePreSignedUrl(ctx, posterImageKey)
	if err != nil {
		return BookResponse{}, err
	}
	book.SetPosterImageLink(url)

	if err = c.Repository.Create(ctx, &book); err != nil {
		return BookResponse{}, err
	}

	return NewBookResponse(book), nil
}

func (c *Catalog) UpdateBook(ctx context.Context, request UpdateBook) error {
	if err := c.Validator.Validate(request); err != nil {
		return err
	}

	existing, err := c.Repository.FindByID(ctx, request.ID)
	if err != nil {
		return err
	}

	updated := request.Update(existing)
	return c.Repository.Update(ctx, &updated)
}

func (c *Catalog) DeleteBook(ctx context.Context, id string) error {
	return c.Repository.Delete(ctx, id)
}
