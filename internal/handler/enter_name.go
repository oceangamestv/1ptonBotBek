package handler

import "github.com/kbgod/illuminate/router"

func (h *Handler) OnName(ctx *router.Context) error {
	return ctx.ReplyVoid("Your name is " + ctx.Message().Text + "!")
}

func (h *Handler) InsideEnterName(ctx *router.Context) error {
	return ctx.ReplyVoid("You are inside enter_name state")
}

func (h *Handler) ExitFromEnterName(ctx *router.Context) error {
	if err := h.svc.RemoveUserBotState(getUserFromContext(ctx.Context)); err != nil {
		return err
	}
	return ctx.ReplyVoid("You are exit from enter_name state")

	// OR
	// return h.Start(ctx)
}

func (h *Handler) SetUserName(ctx *router.Context) error {
	if err := h.svc.SetUserBotState(getUserFromContext(ctx.Context), "enter_name"); err != nil {
		return err
	}
	return ctx.ReplyVoid("Enter your name, or send /cancel")
}
