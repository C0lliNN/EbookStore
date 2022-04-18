//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
)

func InitializeServer() *http.Server {
	wire.Build(config.Provider, auth.Provider, catalog.Provider, shop.Provider, api.Provider)

	return &http.Server{}
}
