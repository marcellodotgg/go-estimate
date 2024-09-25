package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type breakoutController struct {
	pageObject
}

func NewBreakoutController() breakoutController {
	return breakoutController{}
}

func (c breakoutController) Create(ctx *gin.Context) {
	c.reset(ctx)
	ctx.Header("HX-Redirect", "/breakout?id=1238123981273")
	ctx.HTML(http.StatusOK, "breakout/index", c)
}