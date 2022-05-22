package main

import "github.com/c0llinn/ebook-store/internal/config"

func main() {
	config.LoadConfiguration()

	server := NewServer()

	panic(server.Start())
}
