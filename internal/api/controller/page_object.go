package controller

import (
	"os"

	"github.com/gin-gonic/gin"
)

type pageObject struct {
	Hash string
}

func (p *pageObject) reset(_ *gin.Context) {
	if os.Getenv("GIN_MODE") == "release" {
		p.Hash = os.Getenv("BUILD_HASH")
	} else {
		p.Hash = "local"
	}
}
