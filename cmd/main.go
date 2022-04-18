package main

import (
	config2 "github.com/c0llinn/ebook-store/internal/config"
)

func init() {
	config2.InitConfiguration()
	config2.InitLogger()
	config2.LoadMigrations("file:../migrations")
}

func main() {
	server := CreateWebServer()

	err := server.ListenAndServe()
	if err != nil {
		config2.Logger.Fatalf("Could not start web api: %v", err)
	}
}
