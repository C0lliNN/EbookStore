package main

import (
	"github.com/c0llinn/ebook-store/config"
)

func init() {
	config.InitConfiguration()
	config.InitLogger()
	config.LoadMigrations("file:../migration")
}

func main() {
	server := CreateWebServer()

	err := server.ListenAndServe()
	if err != nil {
		config.Logger.Fatalf("Could not start web api: %v", err)
	}
}
