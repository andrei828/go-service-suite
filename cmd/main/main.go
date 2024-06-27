package main

import (
	"github.com/andrei828/go-service-suite/pkg/video"
	"github.com/andrei828/go-service-suite/pkg/webserver"
	"log"
)

func main() {
	// TODO: Implement DI for al the services below
	logger := log.Default()
	uploader := video.NewUploader(logger)
	routeManager := webserver.NewRouteManager(uploader, logger)

	// Kick off webserver
	webServer := webserver.CreateGinWebServer(routeManager, logger)
	if err := webServer.Initialize(); err != nil {
		panic("Failed to initialize web server")
	}

	if err := webServer.Listen(":8080"); err != nil {
		panic("oh no")
	}
}
