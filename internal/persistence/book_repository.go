package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindByQuery(ctx context.Context, query catalog.BookQuery) (paginated catalog.PaginatedBooks, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	db := r.db.WithContext(ctx)

	conditions := r.createConditionsFromCriteria(query.CreateCriteria())

	result := db.Limit(query.Limit).Offset(query.Offset).Where(conditions).Find(&paginated.Books)
	if err = result.Error; err != nil {
		return
	}

	var count int64
	if err = db.Model(&catalog.Book{}).Where(conditions).Count(&count).Error; err != nil {
		return
	}

	paginated.Limit = query.Limit
	paginated.Offset = query.Offset
	paginated.TotalBooks = count
	return
}

func (r *BookRepository) createConditionsFromCriteria(criteria []catalog.Criteria) string {
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

func (r *BookRepository) FindByID(ctx context.Context, id string) (book catalog.Book, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).First(&book, "id = ?", id)
	if err = result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = &ErrEntityNotFound{entity: "book"}
		}
	}

	return
}

func (r *BookRepository) Create(ctx context.Context, book *catalog.Book) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Create(book)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to create a book: %v", err)
		return err
	}

	return nil
}

func (r *BookRepository) Update(ctx context.Context, book *catalog.Book) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Save(book)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&catalog.Book{}, "id = ?", id)
	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected <= 0 {
		return &ErrEntityNotFound{entity: "book"}
	}

	return nil
}
