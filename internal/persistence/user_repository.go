package persistence

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *auth.User) error {
	result := r.db.Create(user)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to save user: %v", err)

		if isConstraintViolationError(err) {
			return &ErrDuplicateKey{key: "email"}
		}

		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (user auth.User, err error) {
	result := r.db.First(&user, "email = ?", email)
	if err = result.Error; err != nil {
		log.Default().Errorf("error trying to find user by email %s: %v", email, err)

		err = &ErrEntityNotFound{entity: "User"}
	}

	return
}

func (r *UserRepository) Update(ctx context.Context, user *auth.User) error {
	result := r.db.Updates(user).Where("id = ?", user.ID)
	if err := result.Error; err != nil {
		log.Default().Errorf("error trying to update user: %v", err)
		return err
	}

	return nil
}

func isConstraintViolationError(err error) bool {
	parsed, ok := err.(*pgconn.PgError)
	return ok && parsed.Code == "23505"
}
