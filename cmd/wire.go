//go:build wireinject
// +build wireinject

package main

import (
	"github.com/c0llinn/ebook-store/config/di"
	"github.com/google/wire"
	"net/http"
)

func CreateWebServer() *http.Server {
	wire.Build(di.Container)
	return &http.Server{}
}
