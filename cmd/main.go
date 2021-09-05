package main

import (
	"fmt"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/env"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/google/uuid"
)

func init() {
	env.InitConfiguration()
	log.InitLogger()
	db.LoadMigrations()
}

func main() {
	repo := SetupApplication()

	user := auth.User{
		ID:        uuid.NewString(),
		FirstName: "Raphael",
		LastName:  "Collin",
		Role:      auth.Admin,
		Email:     "raphael@test.com",
		Password:  "some-password",
	}

	err := repo.Save(&user)
	fmt.Println(err)
}
