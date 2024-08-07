package video

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type EventStreamRequest struct {
	Message string `json:"message" binding:"required"`
}
type Route struct {
	path        string
	method      string
	handler     func(ctx *gin.Context)
	handlerChan func(messagesChan chan string) gin.HandlerFunc
}

func NewSSERoute(method string, path string, handlerChan func(messagesChan chan string) gin.HandlerFunc) *Route {
	return &Route{
		path:        path,
		method:      method,
		handlerChan: handlerChan,
	}
}

func NewRoute(method string, path string, handler func(ctx *gin.Context)) *Route {
	return &Route{
		path:    path,
		method:  method,
		handler: handler,
	}
}

type EventHandler struct {
	routes       []*Route
	logger       *log.Logger
	messagesChan chan string
}

func NewEventHandler(logger *log.Logger) *EventHandler {
	messagesChan := make(chan string)
	return &EventHandler{
		logger:       logger,
		messagesChan: messagesChan,
		routes: []*Route{
			NewRoute(http.MethodPost, "/play_event", playEvent),
			NewRoute(http.MethodPost, "/pause_event", pauseEvent),
			NewSSERoute(http.MethodPost, "/stream_event", streamEvent),
			NewSSERoute(http.MethodGet, "/receive_event", receiveEvent),
		},
	}
}

func (eh *EventHandler) RegisterRoutes(engine *gin.Engine) error {
	for _, route := range eh.routes {
		if route.handler != nil {
			engine.Handle(route.method, route.path, route.handler)
		}
		if route.handlerChan != nil {
			engine.Handle(route.method, route.path, route.handlerChan(eh.messagesChan))
		}
	}
	return nil
}

func streamEvent(messagesChan chan string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request EventStreamRequest
		if err := ctx.ShouldBind(&request); err != nil {
			errorMessage := fmt.Sprintf("request validation error: %s", err.Error())
			BadRequestResponse(ctx, errors.New(errorMessage))
			return
		}

		messagesChan <- request.Message
		CreatedResponse(ctx, &request.Message)

		return
	}
}

func receiveEvent(messagesChan chan string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-messagesChan; ok {
				ctx.SSEvent("message", msg)
				return true
			}
			return false
		})

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
