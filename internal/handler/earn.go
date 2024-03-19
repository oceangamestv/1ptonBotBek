package handler

import (
	"github.com/kbgod/illuminate/router"
)

func (h *Handler) EarnPage(ctx *router.Context) error {
	return ctx.ReplyVoid("Earn page")
}
