package catalog

import (
	"context"
	"fmt"

	"github.com/ebookstore/internal/core/query"
	"github.com/ebookstore/internal/log"
)

type Repository interface {
	FindByQuery(ctx context.Context, query query.Query, page query.Page) (paginated PaginatedBooks, err error)
	FindByID(ctx context.Context, id string) (book Book, err error)
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}

type StorageClient interface {
	GenerateGetPreSignedUrl(ctx context.Context, key string) (string, error)
	GeneratePutPreSignedUrl(ctx context.Context, key string) (string, error)
}

type IDGenerator interface {
	NewID() string
}

type Validator interface {
	Validate(i interface{}) error
}

type Config struct {
	Repository    Repository
	StorageClient StorageClient
	IDGenerator   IDGenerator
	Validator     Validator
}

type Catalog struct {
	Config
}

func New(c Config) *Catalog {
	return &Catalog{Config: c}
}

func (c *Catalog) FindBooks(ctx context.Context, request SearchBooks) (PaginatedBooksResponse, error) {
	log.Infof(ctx, "new request for fetching books")

	paginatedBooks, err := c.Repository.FindByQuery(ctx, request.CreateQuery(), request.CreatePage())
	if err != nil {
		return PaginatedBooksResponse{}, fmt.Errorf("(FindBooks) failed finding books: %w", err)
	}

	imageLinksByBookId := make(map[string][]string)

	for _, book := range paginatedBooks.Books {
		imageLinks, err := c.getPresignedUrlsForBook(ctx, book)
		if err != nil {
			return PaginatedBooksResponse{}, fmt.Errorf("(FindBooks) failed getting urls: %w", err)
		}
		imageLinksByBookId[book.ID] = imageLinks
	}

	return NewPaginatedBooksResponse(paginatedBooks, imageLinksByBookId), nil
}

func (c *Catalog) FindBookByID(ctx context.Context, id string) (BookResponse, error) {
	log.Infof(ctx, "new request for fetching book %s", id)

	book, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		return BookResponse{}, fmt.Errorf("(FindBookByID) failed finding book %s: %w", id, err)
	}

	imageLinks, err := c.getPresignedUrlsForBook(ctx, book)
	if err != nil {
		return BookResponse{}, fmt.Errorf("(FindBookByID) failed getting urls: %w", err)
	}

	return NewBookResponse(book, imageLinks), nil
}

func (c *Catalog) getPresignedUrlsForBook(ctx context.Context, book Book) ([]string, error) {
	imageLinks := make([]string, 0, len(book.Images))
	for _, img := range book.Images {
		url, err := c.StorageClient.GenerateGetPreSignedUrl(ctx, img.ID)
		if err != nil {
			return nil, fmt.Errorf("(getPresignedUrlsForBook) failed generating presigned url: %w", err)
		}
		imageLinks = append(imageLinks, url)
	}

	return imageLinks, nil
}

func (c *Catalog) GetBookContentURL(ctx context.Context, id string) (string, error) {
	log.Infof(ctx, "new request for generating book content url %s", id)

	book, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("(GetBookContent) failed finding book %s: %w", id, err)
	}

	contentUrl, err := c.StorageClient.GenerateGetPreSignedUrl(ctx, book.ContentID)
	if err != nil {
		return "", fmt.Errorf("(GetBookContent) failed generating presigned url: %w", err)
	}

	return contentUrl, nil
}

func (c *Catalog) CreateBook(ctx context.Context, request CreateBook) (BookResponse, error) {
	log.Infof(ctx, "new request for creating book")

	if !isAdmin(ctx) {
		return BookResponse{}, fmt.Errorf("(CreateBook) failed validating access conditions: %w", ErrForbiddenCatalogAccess)
	}

	if err := c.Validator.Validate(request); err != nil {
		return BookResponse{}, fmt.Errorf("(CreateBook) failed validating request: %w", err)
	}

	book := request.Book(c.IDGenerator.NewID())
	if err := c.Repository.Create(ctx, &book); err != nil {
		return BookResponse{}, fmt.Errorf("(CreateBook) failed creating book: %w", err)
	}

	return c.FindBookByID(ctx, book.ID)
}

func (c *Catalog) UpdateBook(ctx context.Context, request UpdateBook) error {
	log.Infof(ctx, "new request for updating book with id %s", request.ID)

	if !isAdmin(ctx) {
		return fmt.Errorf("(UpdateBook) failed validating access conditions: %w", ErrForbiddenCatalogAccess)
	}

	if err := c.Validator.Validate(request); err != nil {
		return fmt.Errorf("(UpdateBook) failed validating request: %w", err)
	}

	existing, err := c.Repository.FindByID(ctx, request.ID)
	if err != nil {
		return fmt.Errorf("(UpdateBook) failed finding book %s: %w", request.ID, err)
	}

	updated := request.Update(existing)
	if err = c.Repository.Update(ctx, &updated); err != nil {
		return fmt.Errorf("(UpdateBook) failed updating book %s: %w", request.ID, err)
	}
	return nil
}

func (c *Catalog) DeleteBook(ctx context.Context, id string) error {
	log.Infof(ctx, "new request for deleting book with id %s", id)

	if !isAdmin(ctx) {
		return fmt.Errorf("(DeleteBook) failed validating access conditions: %w", ErrForbiddenCatalogAccess)
	}

	if err := c.Repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("(DeleteBook) failed deleting book: %w", err)
	}

	return nil
}

func (c *Catalog) GeneratePutPreSignedUrl(ctx context.Context) (PresignURLResponse, error) {
	idGenerator := c.IDGenerator.NewID()
	url, err := c.StorageClient.GeneratePutPreSignedUrl(ctx, idGenerator)
	if err != nil {
		return PresignURLResponse{}, fmt.Errorf("(GeneratePutPreSignedUrl) failed generating presigned url: %w", err)
	}

	return PresignURLResponse{ID: idGenerator, URL: url}, nil
}

func isAdmin(ctx context.Context) bool {
	admin, ok := ctx.Value("admin").(bool)
	if !ok {
		return false
	}

	return admin
}
