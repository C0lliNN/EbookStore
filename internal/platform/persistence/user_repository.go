package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ebookstore/internal/core/auth"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *auth.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Create(user)
	if err := result.Error; err != nil {
		if isConstraintViolationError(err) {
			return &ErrDuplicateKey{key: "email"}
		}

		return fmt.Errorf("(Save) failed running insert statement: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	user := auth.User{}
	result := r.db.WithContext(ctx).First(&user, "email = ?", email)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = &ErrEntityNotFound{entity: "User"}
		}

		return auth.User{}, fmt.Errorf("(FindByEmail) failed executing select query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *auth.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Save(user)
	if err := result.Error; err != nil {
		return fmt.Errorf("(Update) failed running update statement: %w", err)
	}

	return nil
}

func isConstraintViolationError(err error) bool {
	parsed, ok := err.(*pq.Error)
	return ok && parsed.Code == "23505"
}
