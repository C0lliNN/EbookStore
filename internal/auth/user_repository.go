package auth

import (
	"github.com/c0llinn/ebook-store/config/log"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Save(user *User) error {
	result := r.db.Create(user)
	if err := result.Error; err != nil {
		log.Logger.Errorf("error trying to save user: %v", err)
		return err
	}

	return nil
}
