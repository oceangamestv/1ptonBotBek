package handler

import (
	"fmt"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
)

func (h *Handler) EnterAdminPanel(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.Role != entity.UserRoleAdmin {
		return nil
	}

	if err := h.svc.SetUserBotState(u, "admin_panel"); err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}
	return h.AdminMainMenu(ctx)
}

func (h *Handler) AdminMainMenu(ctx *router.Context) error {
	stats, err := h.svc.GetStatistics()
	if err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}
	txt := "<b>🔐 Admin panel</b> (/exit)\n\n"
	txt += "📈 <b>Statistics</b>\n"
	txt += fmt.Sprintf("Users: <b>%d</b>\n", stats.UsersCount)
	txt += fmt.Sprintf("Active users: <b>%d</b>\n", stats.ActiveUsersCount)
	txt += fmt.Sprintf("Created today: <b>%d</b>\n", stats.CreatedToday)
	txt += fmt.Sprintf("Stopped today: <b>%d</b>\n", stats.StoppedToday)
	txt += fmt.Sprintf("Created by referral: <b>%d</b>\n", stats.CreatedByRef)

	buttons := illuminate.NewInlineMenu().
		Row().CallbackBtn("📨 Sending", "sending")

	return ctx.ReplyWithMenuVoid(txt, buttons)
}

func (h *Handler) EnterSendingPanel(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)

	if err := h.svc.SetUserBotState(u, "sending"); err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}
	return ctx.ReplyVoid("Send me a message to send to all users")
}
