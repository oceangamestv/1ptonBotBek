package handler

import (
	"context"
	"fmt"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/service"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
	"runtime/debug"
	"strings"
)

func (h *Handler) CallbackQueryAutoAnswer(ctx *router.Context) error {
	if err := ctx.Next(); err != nil {
		return err
	}
	if ctx.Update.CallbackQuery != nil {
		_, _ = h.bot.AnswerCallbackQuery(ctx.Update.CallbackQuery.Id, nil)
	}
	return nil
}

func (h *Handler) Recovery(ctx *router.Context) error {
	defer func() {
		if err := recover(); err != nil {
			h.svc.Observer.Logger.Error().Interface("panic", err).Msg("fatal error")
			debug.PrintStack()
		}
	}()
	if h.svc.CFG.Debug {
		h.svc.Observer.Logger.Debug().Interface("upd", ctx.Update).Msg("new update")
	}
	return ctx.Next()
}

func (h *Handler) ErrorHandler(ctx *router.Context) error {
	if err := ctx.Next(); err != nil {
		h.svc.Observer.Logger.Error().Interface("upd", ctx.Update).Err(err).Msg("handle update")
	}

	return nil
}

const userCtxKey = "user"

func (h *Handler) UserMiddleware(ctx *router.Context) error {
	tgUser := ctx.Sender()
	if tgUser == nil && ctx.Update.ChosenInlineResult != nil {
		tgUser = &ctx.Update.ChosenInlineResult.From
	}
	if tgUser == nil {
		return fmt.Errorf("sender is nil in context")
	}
	if ctx.Chat() != nil && ctx.Chat().Type != "channel" && ctx.Update.ChatMember != nil {
		return nil
	}
	var isPrivate bool
	if ctx.Chat() != nil && ctx.Chat().Type == "private" {
		isPrivate = true
	}
	var promo *string
	if isCommand(ctx.Update, "/start") {
		args := getCommandArguments(ctx.Update)
		if args != "" {
			promo = new(string)
			*promo = args
		}
	}

	getUserOpts := &service.GetUserOptions{
		TgUser:    tgUser,
		IsPrivate: isPrivate,
		Promo:     promo,
	}

	user, err := h.svc.GetUser(getUserOpts)
	if err != nil {
		return fmt.Errorf("user middleware: %w", err)
	}

	ctx.Context = context.WithValue(ctx.Context, userCtxKey, user)

	if user.BotState != nil {
		ctx.SetState(*user.BotState)
	}
	ctx.SetParseMode(illuminate.ParseModeHTML)
	return ctx.Next()
}

func getUserFromContext(ctx context.Context) *entity.User {
	u, ok := ctx.Value(userCtxKey).(*entity.User)
	if !ok {
		panic("no user in context")
	}

	return u
}

func isCommand(update *illuminate.Update, cmd ...string) bool {
	if update.Message == nil {
		return false
	}
	if update.Message.Entities == nil || len(update.Message.Entities) == 0 {
		return false
	}

	e := update.Message.Entities[0]

	if len(cmd) > 0 {
		return e.Offset == 0 && e.Type == "bot_command" && strings.HasPrefix(update.Message.Text, cmd[0])
	}

	return e.Offset == 0 && e.Type == "bot_command"
}

func getCommandArguments(update *illuminate.Update) string {
	if !isCommand(update) {
		return ""
	}
	e := update.Message.Entities[0]
	if len(update.Message.Text) == int(e.Length) {
		return ""
	}

	return update.Message.Text[e.Length+1:]
}
