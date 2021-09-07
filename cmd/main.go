package main

import (
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/env"
	"github.com/c0llinn/ebook-store/config/log"
)

func init() {
	env.InitConfiguration()
	log.InitLogger()
	db.LoadMigrations()
}

func main() {
	server := CreateWebServer()

	err := server.ListenAndServe()
	if err != nil {
		log.Logger.Fatalf("Could not start web api: %v", err)
	}
}
