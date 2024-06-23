package webserver

import (
	"github.com/gin-gonic/gin"
)

type RouteManager struct {
	server GinWebServer
}

func NewRouteManager() *RouteManager {
	return &RouteManager{}
}

func (rh *RouteManager) RegisterRoutes(engine *gin.Engine) error {
	engine.GET("/public")
}
