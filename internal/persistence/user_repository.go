package persistence

import (
	"context"
	"errors"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
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

		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (user auth.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).First(&user, "email = ?", email)
	if err = result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = &ErrEntityNotFound{entity: "User"}
		}
	}

	return
}

func (r *UserRepository) Update(ctx context.Context, user *auth.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	result := r.db.WithContext(ctx).Save(user)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func isConstraintViolationError(err error) bool {
	parsed, ok := err.(*pq.Error)
	return ok && parsed.Code == "23505"
}
