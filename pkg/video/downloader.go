package video

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Downloader struct {
	handler func(ctx *gin.Context)
	logger  *log.Logger
	route   string
}

func NewDownloader(logger *log.Logger) *Downloader {
	return &Downloader{
		handler: videoDownloaderHandler,
		logger:  logger,
		route:   "/download_video",
	}
}

func (v *Downloader) RegisterRoute(engine *gin.Engine) error {
	engine.GET(v.route, v.handler)
	return nil
}

func videoDownloaderHandler(ctx *gin.Context) {
	ctx.Redirect(http.StatusPermanentRedirect, "/static/bunny/sample_200.mpd")
}
