package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Uploader struct {
	handler func(ctx *gin.Context)
	logger  *log.Logger
	route   string
}

func NewUploader(logger *log.Logger) *Uploader {
	return &Uploader{
		handler: videoUploadHandler,
		logger:  logger,
		route:   "/upload_video",
	}
}

func (v *Uploader) RegisterRoute(engine *gin.Engine) error {
	//engine.MaxMultipartMemory = 8 << 20 // 8 MiB
	engine.POST(v.route, v.handler)
	return nil
}

func videoUploadHandler(ctx *gin.Context) {
	// Multipart form
	form, _ := ctx.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)

		// Upload the file to specific dst.
		err := ctx.SaveUploadedFile(file, "/Users/Andrei/Documents/Go/go-service-suite/iMovie Library.imovielibrary")
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store the files"})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf("%d files uploaded!", len(files))})
}
