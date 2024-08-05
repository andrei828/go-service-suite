package webserver

import (
	"github.com/andrei828/go-service-suite/pkg/video"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type RouteHandler interface {
	RegisterRoute(*gin.Engine) error
}

type RouteManager struct {
	logger       *log.Logger
	uploader     *video.Uploader
	downloader   *video.Downloader
	eventHandler *video.EventHandler
}

func NewRouteManager(
	uploader *video.Uploader,
	downloader *video.Downloader,
	eventHandler *video.EventHandler,
	logger *log.Logger) *RouteManager {
	return &RouteManager{
		logger:       logger,
		uploader:     uploader,
		downloader:   downloader,
		eventHandler: eventHandler,
	}
}

func (rm *RouteManager) RegisterRoutes(engine *gin.Engine) error {
	// register video uploader
	if err := rm.uploader.RegisterRoute(engine); err != nil {
		rm.logger.Fatalf("failed to register the video uploader route.")
		return err
	}

	if err := rm.downloader.RegisterRoute(engine); err != nil {
		rm.logger.Fatalf("failed to register the video downloader route.")
		return err
	}

	if err := rm.eventHandler.RegisterRoutes(engine); err != nil {
		rm.logger.Fatalf("failed to register the video event routes.")
		return err
	}

	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/")
	})

	engine.StaticFS("/static", http.Dir("./static"))
	return nil
}
