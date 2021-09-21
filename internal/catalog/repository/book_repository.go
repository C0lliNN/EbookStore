package repository

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return BookRepository{db: db}
}

func (r BookRepository) FindByQuery(query model.BookQuery) (paginated model.PaginatedBooks, err error) {
	return
}

func (r BookRepository) FindById(id string) (book model.Book, err error) {
	return
}

func (r BookRepository) Create(book *model.Book) error {
	return nil
}

func (r BookRepository) Update(book *model.Book) error {
	return nil
}

func (r BookRepository) Delete(id string) error {
	return nil
}
