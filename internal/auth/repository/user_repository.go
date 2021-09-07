package repository

import (
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth/model"
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
		return err
	}

	return nil
}
