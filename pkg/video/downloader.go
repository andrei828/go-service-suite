package video

import (
	"github.com/gin-gonic/gin"
	"log"
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
	engine.POST(v.route, v.handler)
	return nil
}

func videoDownloaderHandler(ctx *gin.Context) {
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename=movie.mpd")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File("")
}
