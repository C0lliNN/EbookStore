package main

import (
	"fmt"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/env"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/test/factory"
)

func init() {
	env.InitConfiguration()
	log.InitLogger()
	db.LoadMigrations()
}

func main() {
	repo := SetupApplication()

	user := factory.NewUser()
	err := repo.Save(&user)

	fmt.Println(err)
}
