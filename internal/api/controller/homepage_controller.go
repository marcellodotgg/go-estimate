package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type homepageController struct {
	pageObject
}

func NewHomePageController() homepageController {
	return homepageController{}
}

func (c homepageController) HomePage(ctx *gin.Context) {
	c.reset(ctx)
	ctx.HTML(http.StatusOK, "homepage/index", c)
}
