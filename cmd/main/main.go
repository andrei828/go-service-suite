package main

import (
	"github.com/andrei828/go-service-suite/pkg/webserver"
	"log"
)

func main() {
	if err := webserver.CreateGinWebServer(webserver.NewRouteManager(), log.Default()).Listen(":8080"); err != nil {
		panic("oh no")
	}
}
