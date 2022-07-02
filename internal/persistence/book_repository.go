package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ebookstore/internal/catalog"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindByQuery(ctx context.Context, query catalog.BookQuery) (catalog.PaginatedBooks, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	db := r.db.WithContext(ctx)

	conditions := r.createConditionsFromCriteria(query.CreateCriteria())

	paginated := catalog.PaginatedBooks{}
	result := db.Limit(query.Limit).Offset(query.Offset).Where(conditions).Find(&paginated.Books)
	if err := result.Error; err != nil {
		return catalog.PaginatedBooks{}, fmt.Errorf("(FindByQuery) failed running select query: %w", err)
	}

	var count int64
	if err := db.Model(&catalog.Book{}).Where(conditions).Count(&count).Error; err != nil {
		return catalog.PaginatedBooks{}, fmt.Errorf("(FindByQuery) failed running count query: %w", err)
	}

	paginated.Limit = query.Limit
	paginated.Offset = query.Offset
	paginated.TotalBooks = count
	return paginated, nil
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

func (r *BookRepository) FindByID(ctx context.Context, id string) (catalog.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	book := catalog.Book{}
	result := r.db.WithContext(ctx).First(&book, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = &ErrEntityNotFound{entity: "book"}
		}

		return catalog.Book{}, fmt.Errorf("(FindByID) failed running select query: %w", err)
	}

	return book, nil
}

func (r *BookRepository) Create(ctx context.Context, book *catalog.Book) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Create(book)
	if err := result.Error; err != nil {
		return fmt.Errorf("(Create) failed running insert statement: %w", err)
	}

	return nil
}

func (r *BookRepository) Update(ctx context.Context, book *catalog.Book) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Save(book)
	if err := result.Error; err != nil {
		return fmt.Errorf("(Update) failed running update statement: %w", err)
	}

	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&catalog.Book{}, "id = ?", id)
	if err := result.Error; err != nil {
		return fmt.Errorf("(Delete) failed running delete statement: %w", err)
	}

	if result.RowsAffected <= 0 {
		return fmt.Errorf("(Delete) no rows affected: %w", &ErrEntityNotFound{entity: "book"})
	}

	return nil
}
