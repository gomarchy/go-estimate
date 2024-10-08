package controller

import (
	"os"

	"github.com/gin-gonic/gin"
)

type pageObject struct {
	Hash      string
	UserID    string
	ModalType string
}

func (p *pageObject) reset(ctx *gin.Context) {
	if os.Getenv("GIN_MODE") == "release" {
		p.Hash = os.Getenv("BUILD_HASH")
	} else {
		p.Hash = "local"
	}

	p.UserID = ctx.MustGet("user_id").(string)
}
