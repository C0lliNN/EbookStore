package main

import (
	"fmt"
	"github.com/c0llinn/ebook-store/cmd/app"
)

func main() {
	foo := app.SetupApplication()
	fmt.Println("I'm working correctly!", foo)
}
