package main

import "github.com/c0llinn/ebook-store/internal/config"

func main() {
	config.InitConfiguration()

	server := NewServer()

	panic(server.Start())
}
