package main

import (
	"fmt"
	"github.com/c0llinn/ebook-store/cmd/app"
	_ "github.com/c0llinn/ebook-store/config"
	"time"
)

func main() {
	foo := app.SetupApplication()
	fmt.Println("I'm working correctly!", foo)

	time.Sleep(time.Second * 3)
}
