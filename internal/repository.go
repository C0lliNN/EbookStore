package internal

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) HealthTest() {
	row := r.db.Select("SELECT 1").Row()
	var test int
	row.Scan(&test)

	fmt.Println(1)
}
