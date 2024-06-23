package webserver

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type RouteHandler struct {
	name    string
	method  string
	handler *func()
}

type WebServer interface {
	Run()
	RegisterRouteHandler(routeHandler RouteHandler) error
}

type GinWebServer struct {
	internal      *gin.Engine
	routeHandlers map[string]*func()
}

func CreateGinWebServer() *GinWebServer {
	return &GinWebServer{
		internal:      gin.Default(),
		routeHandlers: make(map[string]*func()),
	}
}

func (webServer *GinWebServer) RegisterRouteHandler(routeHandler *RouteHandler) error {
	if routeHandler == nil {
		return errors.New("cannot register a nil routeHandler")
	}

	if _, exists := webServer.routeHandlers[routeHandler.name]; exists {
		return errors.New("route already has a handler")
	}

	webServer.routeHandlers[routeHandler.name] = routeHandler.handler
	webServer.internal.GET(routeHandler.name)
	return nil
}
