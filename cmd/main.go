package main

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth"
)

func main() {
	repo := SetupApplication()

	user := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Password:  "some-password",
	}

	err := repo.Save(&user)
	fmt.Println(err)
}
