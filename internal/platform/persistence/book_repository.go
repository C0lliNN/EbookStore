package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/query"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindByQuery(ctx context.Context, query query.Query, page query.Page) (catalog.PaginatedBooks, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	db := r.db.WithContext(ctx)

	conditions := parseQuery(query)

	paginated := catalog.PaginatedBooks{}
	result := db.Limit(page.Size).Offset(page.Offset()).Where(conditions).Find(&paginated.Books)
	if err := result.Error; err != nil {
		return catalog.PaginatedBooks{}, fmt.Errorf("(FindByQuery) failed running select query: %w", err)
	}

	var count int64
	if err := db.Preload("Images").Model(&catalog.Book{}).Where(conditions).Count(&count).Error; err != nil {
		return catalog.PaginatedBooks{}, fmt.Errorf("(FindByQuery) failed running count query: %w", err)
	}

	paginated.Limit = page.Size
	paginated.Offset = page.Offset()
	paginated.TotalBooks = count
	return paginated, nil
}

func (r *BookRepository) FindByID(ctx context.Context, id string) (catalog.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	book := catalog.Book{}
	result := r.db.WithContext(ctx).Preload("Images").First(&book, "id = ?", id)
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

	if err := r.db.Model(&book).Association("Images").Replace(book.Images); err != nil {
		return fmt.Errorf("(Update) failed running update statement: %w", err)
	}

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
