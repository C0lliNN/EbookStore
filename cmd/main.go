package main

import "github.com/c0llinn/ebook-store/internal/config"

func init() {
	config.InitConfiguration()
	config.LoadMigrations("file:../migrations")
}

func main() {
	server := NewServer()

	panic(server.Start())
}
