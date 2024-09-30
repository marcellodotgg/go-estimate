package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/infrastructure/websocket"
)

type websocketController struct{}

func NewWebSocketController() websocketController {
	return websocketController{}
}

func (c websocketController) Handle(ctx *gin.Context) {
	channel := ctx.Param("channel")
	websocket.Manager.HandleRequestWithKeys(
		ctx.Writer,
		ctx.Request,
		map[string]interface{}{
			"channel": channel,
		},
	)
}
