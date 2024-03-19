package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
	"strings"
	"time"
)

func (h *Handler) PrepareMessage(ctx *router.Context) error {
	//menu := illuminate.NewInlineMenu()
	//menu.Row().URLBtn("Google", "https://google.com")
	//_, err := h.bot.CopyMessage(ctx.ChatID(), ctx.ChatID(), ctx.Update.Message.MessageID, &illuminate.CopyMessageOpts{
	//	ReplyMarkup: menu,
	//})
	if err := h.svc.SetUserBotState(
		getUserFromContext(ctx.Context), "sending", NewSendingMessage(ctx.Update.Message.MessageID),
	); err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}

	return ctx.ReplyVoid("Now, send me a keyboard in the format: \n\n" +
		"<code>text - https://link1\n" +
		"text2 - https://link2</code>" +
		"\n\nOr /skip to send without a keyboard")
}

type SendingMessageKeyboard struct {
	Text string `json:"text"`
	Link string `json:"link"`
}
type SendingMessage struct {
	MessageID int64                    `json:"message_id"`
	Keyboard  []SendingMessageKeyboard `json:"keyboard"`
}

func NewSendingMessage(messageID int64) string {
	j, _ := json.Marshal(SendingMessage{MessageID: messageID})
	return string(j)
}

func NewSendingMessageFromString(message string) (SendingMessage, error) {
	var sm SendingMessage
	if err := json.Unmarshal([]byte(message), &sm); err != nil {
		return sm, err
	}
	return sm, nil
}

func NewSendingMessageWithKeyboard(messageID int64, keyboard []SendingMessageKeyboard) string {
	j, _ := json.Marshal(SendingMessage{
		MessageID: messageID,
		Keyboard:  keyboard,
	})
	return string(j)
}

func (h *Handler) WithKeyboard(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.BotStateContext == nil {
		return ctx.AnswerAlertVoid("You should send a message first")
	}
	sendingMessage, err := NewSendingMessageFromString(*u.BotStateContext)
	if err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}
	menu := illuminate.NewInlineMenu()
	lines := strings.Split(ctx.Update.Message.Text, "\n")

	keyboard := make([]SendingMessageKeyboard, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, " - ")
		switch {
		case parts[1] == "mine":
			menu.Row().WebAppBtn(parts[0], h.svc.CFG.FrontendURL)
		default:
			menu.Row().URLBtn(parts[0], parts[1])
		}
		keyboard = append(keyboard, SendingMessageKeyboard{
			Text: parts[0],
			Link: parts[1],
		})
	}
	if err := h.svc.SetUserBotState(
		getUserFromContext(ctx.Context), "sending", NewSendingMessageWithKeyboard(sendingMessage.MessageID, keyboard),
	); err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}

	menu.Row().CallbackBtn("📤 Send", "send")

	_, err = h.bot.CopyMessage(ctx.ChatID(), ctx.ChatID(), sendingMessage.MessageID, &illuminate.CopyMessageOpts{
		ReplyMarkup: menu,
	})

	return err
}

func (h *Handler) WithoutKeyboard(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.BotStateContext == nil {
		return ctx.AnswerAlertVoid("You should send a message first")
	}
	sendingMessage, err := NewSendingMessageFromString(*u.BotStateContext)
	if err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}

	menu := illuminate.NewInlineMenu()
	menu.Row().CallbackBtn("📤 Send", "send")
	_, err = h.bot.CopyMessage(ctx.ChatID(), ctx.ChatID(), sendingMessage.MessageID, &illuminate.CopyMessageOpts{
		ReplyMarkup: menu,
	})
	return err
}

func (h *Handler) Send(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.BotStateContext == nil {
		return ctx.AnswerAlertVoid("You should send a message first")
	}
	sendingMessage, err := NewSendingMessageFromString(*u.BotStateContext)
	if err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}

	sendingMessageChan <- SendingMessageTask{
		SendingMessage:   sendingMessage,
		AdminChatID:      ctx.ChatID(),
		ReplyToMessageID: ctx.Message().MessageID,
	}

	_ = ctx.AnswerAlertVoid("Sending started")

	_ = h.svc.RemoveUserBotState(u)

	return h.Start(ctx)

	//opts := &illuminate.CopyMessageOpts{}
	//if len(sendingMessage.Keyboard) > 0 {
	//	menu := illuminate.NewInlineMenu()
	//	for _, button := range sendingMessage.Keyboard {
	//		menu.Row().URLBtn(button.Text, button.Link)
	//	}
	//	opts.ReplyMarkup = menu
	//}
	//
	//var (
	//	success int64
	//	failed  int64
	//)
	//start := time.Now()
	//err = h.svc.BatchActiveUsers(10000, func(user entity.User) {
	//	_, err := h.bot.CopyMessage(user.ID, ctx.ChatID(), sendingMessage.MessageID, opts)
	//	time.Sleep(40 * time.Millisecond)
	//	if err != nil {
	//		failed++
	//		h.svc.Observer.Logger.Error().Err(err).Msg("send message")
	//		return
	//	}
	//	success++
	//})
	//if err != nil {
	//	return ctx.AnswerAlertVoid(err.Error())
	//}
	//
	//return ctx.ReplyVoid(
	//	fmt.Sprintf(
	//		"<b>Sending finished!</b>\n\nSent: <code>%d</code>\nFailed: <code>%d</code>\nTime elapsed: <code>%s</code>",
	//		success, failed, time.Since(start),
	//	), &illuminate.SendMessageOpts{
	//		ReplyParameters: &illuminate.ReplyParameters{
	//			MessageID: ctx.Message().MessageID,
	//		},
	//	},
	//)
}

type SendingMessageTask struct {
	SendingMessage
	AdminChatID      int64
	ReplyToMessageID int64
}

var sendingMessageChan = make(chan SendingMessageTask)

func (h *Handler) StartSenderWorker() {
	go func() {
		for sm := range sendingMessageChan {
			go func(sm SendingMessageTask) {
				opts := &illuminate.CopyMessageOpts{}
				if len(sm.Keyboard) > 0 {
					menu := illuminate.NewInlineMenu()
					for _, button := range sm.Keyboard {
						menu.Row().URLBtn(button.Text, button.Link)
					}
					opts.ReplyMarkup = menu
				}

				var (
					success int64
					failed  int64
				)
				start := time.Now()
				err := h.svc.BatchActiveUsers(10000, func(user entity.User) {
					_, err := h.bot.CopyMessage(user.ID, sm.AdminChatID, sm.MessageID, opts)
					time.Sleep(40 * time.Millisecond)
					if err != nil {
						failed++
						h.svc.Observer.Logger.Error().Err(err).Msg("send message")
						return
					}
					success++
				})
				if err != nil {
					h.svc.Observer.Logger.Error().Err(err).Msg("send message")
					return
				}

				_, _ = h.bot.SendMessage(sm.AdminChatID, fmt.Sprintf(
					"<b>Sending finished!</b>\n\nSent: <code>%d</code>\nFailed: <code>%d</code>\nTime elapsed: <code>%s</code>",
					success, failed, time.Since(start),
				), &illuminate.SendMessageOpts{
					ParseMode: illuminate.ParseModeHTML,
					ReplyParameters: &illuminate.ReplyParameters{
						MessageID: sm.ReplyToMessageID,
					},
				})
			}(sm)
		}
	}()
}

func (h *Handler) ExitFromSending(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if err := h.svc.SetUserBotState(u, "admin_panel"); err != nil {
		return ctx.AnswerAlertVoid(err.Error())
	}
	return h.AdminMainMenu(ctx)
}
