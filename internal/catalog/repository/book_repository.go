package repository

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/config"
	"gorm.io/gorm"
	"strings"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return BookRepository{db: db}
}

func (r BookRepository) FindByQuery(query model.BookQuery) (paginated model.PaginatedBooks, err error) {
	conditions := r.createConditionsFromCriteria(query.CreateCriteria())

	result := r.db.Limit(query.Limit).Offset(query.Offset).Where(conditions).Find(&paginated.Books)
	if err = result.Error; err != nil {
		config.Logger.Error("error trying to find books by query: ", err)
		return
	}

	var count int64
	if err = r.db.Model(&model.Book{}).Where(conditions).Count(&count).Error; err != nil {
		config.Logger.Error("error trying to cound books: ", err)
		return
	}

	paginated.Limit = query.Limit
	paginated.Offset = query.Offset
	paginated.TotalBooks = count
	return
}

func (r BookRepository) createConditionsFromCriteria(criteria []model.Criteria) string {
	conditions := make([]string, 0, len(criteria))
	for _, c := range criteria {
		if !c.IsEmpty() {
			if parsed, ok := c.Value.(string); ok {
				conditions = append(conditions, fmt.Sprintf("%s %s '%s'", c.Field, c.Operator, parsed))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s %s %v", c.Field, c.Operator, c.Value))
			}
		}
	}

	return strings.Join(conditions, " AND ")
}

func (r BookRepository) FindByID(id string) (book model.Book, err error) {
	result := r.db.First(&book, "id = ?", id)
	if err = result.Error; err != nil {
		config.Logger.Errorf("error trying to find book by id %s: %v", id, err)
		err = &common.ErrEntityNotFound{Entity: "Book", Err: err}
	}

	return
}

func (r BookRepository) Create(book *model.Book) error {
	result := r.db.Create(book)
	if err := result.Error; err != nil {
		config.Logger.Error("error trying to create a book: ", err)
		return err
	}

	return nil
}

func (r BookRepository) Update(book *model.Book) error {
	result := r.db.Updates(book).Where("id = ?", book.ID)
	if err := result.Error; err != nil {
		config.Logger.Errorf("error trying to update the book wiht id %s: %v", book.ID, err)
		return err
	}

	return nil
}

func (r BookRepository) Delete(id string) error {
	result := r.db.Delete(&model.Book{}, "id = ?", id)
	if err := result.Error; err != nil {
		config.Logger.Errorf("error trying to delete the book with id %s: %v", id, err)
		return err
	}

	if result.RowsAffected <= 0 {
		return &common.ErrEntityNotFound{Entity: "Book", Err: fmt.Errorf("no rows affected")}
	}

	return nil
}
