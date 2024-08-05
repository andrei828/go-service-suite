package webserver

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type WebServer interface {
	Listen(port string) error
	Initialize() error
}

type GinWebServer struct {
	engine       *gin.Engine
	logger       *log.Logger
	routeManager *RouteManager
}

func CreateGinWebServer(routeManager *RouteManager, logger *log.Logger, opts ...gin.OptionFunc) *GinWebServer {
	return &GinWebServer{
		engine:       gin.New(opts...),
		logger:       logger,
		routeManager: routeManager,
	}
}

func (webServer *GinWebServer) Initialize() error {
	if err := webServer.routeManager.RegisterRoutes(webServer.engine); err != nil {
		return err
	}
	return nil
}

func (webServer *GinWebServer) Listen(port string) error {
	srv := &http.Server{
		Addr:    port,
		Handler: webServer.engine.Handler(),
	}

	var err error
	go func() {
		webServer.logger.Println("Listening...")
		if err = http3.ListenAndServeTLS(":443",
			"/Users/andreiliviu/Documents/Go/go-service-suite/example.com+5.pem",
			"/Users/andreiliviu/Documents/Go/go-service-suite/example.com+5-key.pem",
			webServer.engine.Handler()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			webServer.logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	webServer.logger.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	if shutdownError := srv.Shutdown(ctx); shutdownError != nil {
		webServer.logger.Fatal("Server Shutdown:", shutdownError)
	}

	select {
	case <-ctx.Done():
		webServer.logger.Println("timeout of 5 seconds.")
	}
	webServer.logger.Println("Server exiting")
	return err
}
