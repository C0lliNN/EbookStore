package main

import (
	"github.com/c0llinn/ebook-store/cmd/app"
	_ "github.com/c0llinn/ebook-store/config"
)

func main() {
	repo := app.SetupApplication()

	repo.HealthTest()
}
