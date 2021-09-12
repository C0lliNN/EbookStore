package repository

import (
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Save(user *model.User) error {
	result := r.db.Create(user)
	if err := result.Error; err != nil {
		log.Logger.Error("error trying to save user: ", err)

		if parsed, ok := err.(*pgconn.PgError); ok && parsed.Code == "23505" {
			return &model.ErrDuplicateKey{Key: "email", Err: err}
		}

		return err
	}

	return nil
}

func (r UserRepository) FindByEmail(email string) (user model.User, err error) {
	result := r.db.First(&user, "email = ?", email)
	if err = result.Error; err != nil {
		log.Logger.Errorf("error trying to find user by email %s: %v", email, err)

		err = &model.ErrEntityNotFound{Entity: "User", Err: err}
	}

	return
}
