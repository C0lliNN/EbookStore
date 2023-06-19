package server

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
	Public  bool
}

func (r Route) IsPublic() bool {
	return r.Public
}
