package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/service"
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

	if !c.breakoutService.Exists(c.ID) {
		ctx.HTML(http.StatusNotFound, "404", c)
		return
	}

	ctx.HTML(http.StatusOK, "breakout/index", c)
}

func (c breakoutController) Create(ctx *gin.Context) {
	c.reset(ctx)
	breakout, err := c.breakoutService.Create(ctx.MustGet("user_id").(string))

	if err != nil {
		ctx.HTML(http.StatusNotFound, "404", c)
		return
	}

	ctx.Header("HX-Redirect", fmt.Sprintf("/breakout?id=%s", breakout.ID))
}
