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

	engine.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"result": "success"})
	})

	engine.StaticFS("/static", http.Dir("./static"))
	engine.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	return nil
}
