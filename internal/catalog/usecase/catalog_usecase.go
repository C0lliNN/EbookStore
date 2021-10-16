package usecase

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"io"
)

type Repository interface {
	FindByQuery(query model.BookQuery) (paginated model.PaginatedBooks, err error)
	FindByID(id string) (book model.Book, err error)
	Create(book *model.Book) error
	Update(book *model.Book) error
	Delete(id string) error
}

type StorageClient interface {
	GeneratePreSignedUrl(key string) (string, error)
	SaveFile(key string, contentType string, content io.ReadSeeker) error
	RetrieveFile(key string) (io.ReadCloser, error)
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

func (u CatalogUseCase) FindBooks(query model.BookQuery) (paginatedBooks model.PaginatedBooks, err error) {
	paginatedBooks, err = u.repo.FindByQuery(query)
	if err != nil {
		return
	}

	for i := range paginatedBooks.Books {
		imageKey := paginatedBooks.Books[i].PosterImageBucketKey

		var url string
		url, err = u.storageClient.GeneratePreSignedUrl(imageKey)
		if err != nil {
			return
		}

		paginatedBooks.Books[i].SetPosterImageLink(url)
	}

	return
}

func (u CatalogUseCase) FindBookByID(id string) (book model.Book, err error) {
	book, err = u.repo.FindByID(id)
	if err != nil {
		return
	}

	imageKey := book.PosterImageBucketKey
	url, err := u.storageClient.GeneratePreSignedUrl(imageKey)
	if err != nil {
		return
	}

	book.SetPosterImageLink(url)
	return
}

func (u CatalogUseCase) GetBookContent(id string) (io.ReadCloser, error) {
	book, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return u.storageClient.RetrieveFile(book.ContentBucketKey)
}

func (u CatalogUseCase) CreateBook(book *model.Book, posterImage io.ReadSeeker, bookContent io.ReadSeeker) error {
	posterImageKey := u.filenameGenerator.NewUniqueName("poster_" + book.Title)
	contentKey := u.filenameGenerator.NewUniqueName("content_" + book.Title)

	if err := u.storageClient.SaveFile(posterImageKey, "image/jpeg", posterImage); err != nil {
		return err
	}

	if err := u.storageClient.SaveFile(contentKey, "application/pdf", bookContent); err != nil {
		return err
	}

	book.PosterImageBucketKey = posterImageKey
	book.ContentBucketKey = contentKey

	url, err := u.storageClient.GeneratePreSignedUrl(posterImageKey)
	if err != nil {
		return err
	}
	book.SetPosterImageLink(url)

	return u.repo.Create(book)
}

func (u CatalogUseCase) UpdateBook(book *model.Book) error {
	return u.repo.Update(book)
}

func (u CatalogUseCase) DeleteBook(id string) error {
	return u.repo.Delete(id)
}
