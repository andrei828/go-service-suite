package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type Event struct {
	Message   string `json:"message" binding:"required"`
	SessionId string `json:"sessionId" binding:"required"`
}

type Connection chan Event

type EventHandler struct {
	logger      *log.Logger
	connections map[Connection]bool
}

func NewEventHandler(logger *log.Logger) *EventHandler {
	return &EventHandler{
		logger:      logger,
		connections: map[Connection]bool{},
	}
}

func (eh *EventHandler) RegisterRoutes(engine *gin.Engine) error {
	engine.POST("/play_event", playEvent)
	engine.POST("/pause_event", pauseEvent)
	engine.POST("/register_event", registerEvent(eh))
	engine.GET("/stream_event", headersMiddleware(), registerConn(eh), streamEvent)
	return nil
}

func registerConn(eh *EventHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := make(Connection)
		defer func() {
			close(conn)
			delete(eh.connections, conn)
		}()
		eh.connections[conn] = true
		ctx.Set("conn", conn)
		ctx.Next()
	}
}

func streamEvent(ctx *gin.Context) {
	v, ok := ctx.Get("conn")
	if !ok {
		return
	}
	conn, ok := v.(Connection)
	if !ok {
		return
	}
	ctx.Stream(func(w io.Writer) bool {
		event, ok := <-conn
		if !ok {
			return false
		}
		jsonEvent, err := json.Marshal(event)
		if err != nil {
			return false
		}
		ctx.SSEvent("message", string(jsonEvent))
		return true
	})
}

func registerEvent(eh *EventHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Event
		eh.logger.Println("Connections: %d", len(eh.connections), eh.connections)
		if err := ctx.ShouldBind(&request); err != nil {
			errorMessage := fmt.Sprintf("request validation error: %s", err.Error())
			BadRequestResponse(ctx, errors.New(errorMessage))
			return
		}

		for conn := range eh.connections {
			conn <- request
		}

		CreatedResponse(ctx, &request.Message)
		return
	}
}

func playEvent(ctx *gin.Context) {

}

func pauseEvent(ctx *gin.Context) {

}

type JSendFailResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type JSendSuccessResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data,omitempty"`
}

func BadRequestResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusBadRequest,
		JSendFailResponse[string]{
			Status: "fail",
			Data:   error.Error(),
		},
	)

	return
}

func CreatedResponse[T interface{}](c *gin.Context, i *T) {
	c.JSON(
		http.StatusCreated,
		JSendSuccessResponse[T]{
			Status: "success",
			Data:   *i,
		},
	)

	return
}

func headersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
