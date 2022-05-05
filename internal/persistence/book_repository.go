package persistence

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/log"
	"gorm.io/gorm"
	"strings"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return BookRepository{db: db}
}

func (r BookRepository) FindByQuery(ctx context.Context, query catalog.BookQuery) (paginated catalog.PaginatedBooks, err error) {
	conditions := r.createConditionsFromCriteria(query.CreateCriteria())

	result := r.db.Limit(query.Limit).Offset(query.Offset).Where(conditions).Find(&paginated.Books)
	if err = result.Error; err != nil {
		log.Default().Errorf("error trying to find books by query: %v", err)
		return
	}

	var count int64
	if err = r.db.Model(&catalog.Book{}).Where(conditions).Count(&count).Error; err != nil {
		log.Default().Errorf("error trying to count books: %v", err)
		return
	}

	paginated.Limit = query.Limit
	paginated.Offset = query.Offset
	paginated.TotalBooks = count
	return
}

func (r BookRepository) createConditionsFromCriteria(criteria []catalog.Criteria) string {
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

func (r BookRepository) FindByID(ctx context.Context, id string) (book catalog.Book, err error) {
	result := r.db.First(&book, "id = ?", id)
	if err = result.Error; err != nil {
		log.Default().Errorf("error trying to find book by id %s: %v", id, err)
		err = &ErrEntityNotFound{entity: "book"}
	}

	return
}

func (r BookRepository) Create(ctx context.Context, book *catalog.Book) error {
	result := r.db.Create(book)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to create a book: %v", err)
		return err
	}

	return nil
}

func (r BookRepository) Update(ctx context.Context, book *catalog.Book) error {
	result := r.db.Updates(book).Where("id = ?", book.ID)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to update the book wiht id %s: %v", book.ID, err)
		return err
	}

	return nil
}

func (r BookRepository) Delete(ctx context.Context, id string) error {
	result := r.db.Delete(&catalog.Book{}, "id = ?", id)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to delete the book with id %s: %v", id, err)
		return err
	}

	if result.RowsAffected <= 0 {
		return &ErrEntityNotFound{entity: "book"}
	}

	return nil
}
