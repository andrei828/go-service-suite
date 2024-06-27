package webserver

import (
	"github.com/andrei828/go-service-suite/pkg/video"
	"github.com/gin-gonic/gin"
	"log"
)

type RouteHandler interface {
	RegisterRoute(*gin.Engine) error
}

type RouteManager struct {
	logger   *log.Logger
	uploader *video.Uploader
}

func NewRouteManager(uploader *video.Uploader, logger *log.Logger) *RouteManager {
	return &RouteManager{
		logger:   logger,
		uploader: uploader,
	}
}

func (rm *RouteManager) RegisterRoutes(engine *gin.Engine) error {
	// register video uploader
	if err := rm.uploader.RegisterRoute(engine); err != nil {
		rm.logger.Fatalf("failed to register the video uploader route")
		return err
	}

	return nil
}
