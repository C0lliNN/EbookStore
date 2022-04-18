package main

import "github.com/c0llinn/ebook-store/internal/config"

func init() {
	config.InitConfiguration()
	config.InitLogger()
	config.LoadMigrations("file:../migrations")
}

func main() {
	server := CreateWebServer()

	panic(server.ListenAndServe())
}
