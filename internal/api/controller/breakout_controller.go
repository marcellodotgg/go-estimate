package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomarchy/estimate/internal/domain"
	"github.com/gomarchy/estimate/internal/infrastructure/database"
	"github.com/gomarchy/estimate/internal/service"
)

type breakoutController struct {
	pageObject
	Breakout        domain.Breakout
	CurrentUser     domain.User
	breakoutService service.BreakoutService
	voteService     service.VoteService
}

func NewBreakoutController() breakoutController {
	return breakoutController{
		breakoutService: service.NewBreakoutService(),
		voteService:     service.NewVoteService(),
	}
}

func (c *breakoutController) load(ctx *gin.Context) error {
	breakout, err := c.breakoutService.FindByID(ctx.Param("id"))

	if err != nil {
		return err
	}

	c.Breakout = breakout

	database.DB.First(&c.CurrentUser, "user_id = ? AND breakout_id = ?", c.UserID, breakout.ID)
	return nil
}

func (c breakoutController) Index(ctx *gin.Context) {
	c.reset(ctx)
	err := c.load(ctx)

	if err != nil {
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

	ctx.Header("HX-Redirect", fmt.Sprintf("/breakout/%s", breakout.ID))
}

func (c breakoutController) Vote(ctx *gin.Context) {
	c.reset(ctx)
	err := c.load(ctx)

	if err != nil {
		ctx.HTML(http.StatusNotFound, "404", c)
		return
	}

	c.voteService.Vote(&c.CurrentUser, ctx.Query("value"))

	ctx.HTML(http.StatusOK, "vote/index", c.Breakout)
}

func (c breakoutController) Reset(ctx *gin.Context) {
	c.reset(ctx)
	err := c.load(ctx)

	if err != nil {
		ctx.HTML(http.StatusNotFound, "404", c)
		return
	}

	c.voteService.Reset(&c.Breakout)

	ctx.HTML(http.StatusNoContent, "", nil)
}

func (c breakoutController) ShowVotes(ctx *gin.Context) {
	c.reset(ctx)
	err := c.load(ctx)

	if err != nil {
		ctx.HTML(http.StatusNotFound, "404", c)
		return
	}

	c.voteService.ShowVotes(&c.Breakout)

	ctx.HTML(http.StatusNoContent, "", nil)
}

func (c breakoutController) UpdateDisplayNameModal(ctx *gin.Context) {
	c.reset(ctx)
	c.load(ctx)
	c.ModalType = "update_display_name"

	ctx.HTML(http.StatusOK, "modal", c)
}

func (c breakoutController) UpdateDisplayName(ctx *gin.Context) {
	c.reset(ctx)
	c.load(ctx)

	form := GetForm(ctx)
	c.CurrentUser.Name = form.Data["display_name"]

	if err := c.breakoutService.UpdateUser(c.CurrentUser); err != nil {
		return
	}

	ctx.Header("HX-Trigger", "closeModal, confirmedName")
}
