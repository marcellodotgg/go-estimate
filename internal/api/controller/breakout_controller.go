package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/service"
	"github.com/google/uuid"
)

type breakoutController struct {
	pageObject
	ID              string
	UserID          string
	breakoutService service.BreakoutService
}

func NewBreakoutController() breakoutController {
	return breakoutController{
		breakoutService: service.NewBreakoutService(),
	}
}

func (c breakoutController) Index(ctx *gin.Context) {
	c.reset(ctx)
	c.ID, _ = ctx.GetQuery("id")
	c.UserID = ctx.MustGet("user_id").(string)

	if _, exists := service.Channels[c.ID]; !exists {
		ctx.HTML(http.StatusNotFound, "404", c)
		return
	}

	ctx.HTML(http.StatusOK, "breakout/index", c)
}

func (c breakoutController) Create(ctx *gin.Context) {
	c.reset(ctx)
	breakoutID := uuid.NewString()
	c.breakoutService.Create(breakoutID, ctx.MustGet("user_id").(string))
	ctx.Header("HX-Redirect", fmt.Sprintf("/breakout?id=%s", breakoutID))
}
