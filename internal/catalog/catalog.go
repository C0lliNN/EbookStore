package catalog

import (
	"context"
	"fmt"
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
		return PaginatedBooksResponse{}, fmt.Errorf("FindBooks) failed finding books: %w", err)
	}

	for i := range paginatedBooks.Books {
		imageKey := paginatedBooks.Books[i].PosterImageBucketKey

		var url string
		url, err = c.StorageClient.GeneratePreSignedUrl(ctx, imageKey)
		if err != nil {
			bookId := paginatedBooks.Books[i].ID
			return PaginatedBooksResponse{}, fmt.Errorf("FindBooks] failed generating url for book %s: %w", bookId, err)
		}

		paginatedBooks.Books[i].SetPosterImageLink(url)
	}

	return NewPaginatedBooksResponse(paginatedBooks), nil
}

func (c *Catalog) FindBookByID(ctx context.Context, id string) (BookResponse, error) {
	book, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		return BookResponse{}, fmt.Errorf("FindBookByID) failed finding book %s: %w", id, err)
	}

	imageKey := book.PosterImageBucketKey
	url, err := c.StorageClient.GeneratePreSignedUrl(ctx, imageKey)
	if err != nil {
		return BookResponse{}, fmt.Errorf("FindBookByID) failed generating presigned url: %w", err)
	}

	book.SetPosterImageLink(url)
	return NewBookResponse(book), nil
}

func (c *Catalog) GetBookContent(ctx context.Context, id string) (io.ReadCloser, error) {
	book, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetBookContent) failed finding book %s: %w", id, err)
	}

	content, err := c.StorageClient.RetrieveFile(ctx, book.ContentBucketKey)
	if err != nil {
		return nil, fmt.Errorf("GetBookContent) failed retrieving book content: %w", err)
	}

	return content, nil
}

func (c *Catalog) CreateBook(ctx context.Context, request CreateBook) (BookResponse, error) {
	if !isAdmin(ctx) {
		return BookResponse{}, fmt.Errorf("CreateBook) failed validating access conditions: %w", ErrForbiddenCatalogAccess)
	}

	if err := c.Validator.Validate(request); err != nil {
		return BookResponse{}, fmt.Errorf("CreateBook) failed validating request: %w", err)
	}

	book := request.Book(c.IDGenerator.NewID())
	posterImageKey := c.FilenameGenerator.NewUniqueName("poster_" + book.Title)
	contentKey := c.FilenameGenerator.NewUniqueName("content_" + book.Title)

	if err := c.StorageClient.SaveFile(ctx, posterImageKey, "image/jpeg", request.PosterImage); err != nil {
		return BookResponse{}, fmt.Errorf("CreateBook) failed saving poster: %w", err)
	}

	if err := c.StorageClient.SaveFile(ctx, contentKey, "application/pdf", request.BookContent); err != nil {
		return BookResponse{}, fmt.Errorf("CreateBook) failed saving content: %w", err)
	}

	book.PosterImageBucketKey = posterImageKey
	book.ContentBucketKey = contentKey

	url, err := c.StorageClient.GeneratePreSignedUrl(ctx, posterImageKey)
	if err != nil {
		return BookResponse{}, fmt.Errorf("CreateBook) failed generating url: %w", err)
	}
	book.SetPosterImageLink(url)

	if err = c.Repository.Create(ctx, &book); err != nil {
		return BookResponse{}, fmt.Errorf("CreateBook) failed creating book: %w", err)
	}

	return NewBookResponse(book), nil
}

func (c *Catalog) UpdateBook(ctx context.Context, request UpdateBook) error {
	if !isAdmin(ctx) {
		return fmt.Errorf("UpdateBook) failed validating access conditions: %w", ErrForbiddenCatalogAccess)
	}

	if err := c.Validator.Validate(request); err != nil {
		return fmt.Errorf("UpdateBook) failed validating request: %w", err)
	}

	existing, err := c.Repository.FindByID(ctx, request.ID)
	if err != nil {
		return fmt.Errorf("UpdateBook) failed finding book %s: %w", request.ID, err)
	}

	updated := request.Update(existing)
	if err = c.Repository.Update(ctx, &updated); err != nil {
		return fmt.Errorf("UpdateBook) failed updating book %s: %w", request.ID, err)
	}
	return nil
}

func (c *Catalog) DeleteBook(ctx context.Context, id string) error {
	if !isAdmin(ctx) {
		return fmt.Errorf("DeleteBook) failed validating access conditions: %w", ErrForbiddenCatalogAccess)
	}

	if err := c.Repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteBook) failed deleting book: %w", err)
	}

	return nil
}

func isAdmin(ctx context.Context) bool {
	admin, ok := ctx.Value("admin").(bool)
	if !ok {
		return false
	}

	return admin
}
