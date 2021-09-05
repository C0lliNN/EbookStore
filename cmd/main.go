package main

import (
	_ "github.com/c0llinn/ebook-store/config"
)

func main() {
	repo := SetupApplication()

	repo.HealthTest()
}
