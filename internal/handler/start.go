package handler

import (
	"fmt"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
	"time"
)

var startMenuText = map[string]string{
	"en": "Hello %s, I'm <b>%s</b>\n" +
		"News: %s\n" +
		"Tap on the coin and watch your balance grow.\n\n" +
		"<b>🪙 Your coins:</b> %d\n" +
		"<b>💵 Your balance:</b> <code>$%s</code>",
	"ru": "Привет %s, я <b>%s</b>\n" +
		"Новости: %s\n" +
		"Коснитесь монеты и наблюдайте, как растет ваш баланс.\n\n" +
		"<b>🪙 Ваши монеты:</b> %d\n" +
		"<b>💵 Ваш баланс:</b> <code>$%s</code>",
	"uk": "Привіт %s, я <b>%s</b>\n" +
		"Новини: %s\n" +
		"Торкніться монети і спостерігайте, як зростає ваш баланс.\n\n" +
		"<b>🪙 Монети:</b> %d\n" +
		"<b>💵 Ваш баланс:</b> <code>$%s</code>",
}

var premiumText = map[string]string{
	"en": "🏵 Premium until %s UTC",
	"ru": "🏵 Премиум до %s UTC",
	"uk": "🏵 Преміум до %s UTC",
}

var mineButtonText = map[string]string{
	"en": "🕹 Play",
	"ru": "🕹 Играть",
	"uk": "🕹 Грати",
	"tr": "🕹 Oyna",
	"es": "🕹 Jugar",
	"de": "🕹 Spielen",
	"fr": "🕹 Jouer",
	"it": "🕹 Giocare",
	"pt": "🕹 Jogar",
	"nl": "🕹 Spelen",
	"pl": "🕹 Grać",
	"ro": "🕹 Jucați",
}

var premiumButtonText = map[string]string{
	"en": "🏵 Buy premium",
	"ru": "🏵 Купить премиум",
	"uk": "🏵 Придбати преміум",
}

var referralButtonText = map[string]string{
	"en": "🔗 Referral earning",
	"ru": "🔗 Заработок на рефералах",
	"uk": "🔗 Заробіток на рефералах",
}

func (h *Handler) Start(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)

	menu := illuminate.NewInlineMenu().
		Row().WebAppBtn(langf(mineButtonText, u.LanguageCode), h.svc.CFG.FrontendURL).
		Row().CallbackBtn(langf(premiumButtonText, u.LanguageCode), "premium").
		Row().CallbackBtn(langf(referralButtonText, u.LanguageCode), "referral")

	if u.Role == entity.UserRoleAdmin {
		menu.Row().CallbackBtn("⚙️ Admin panel", "admin_panel")
	}
	txt := langf(startMenuText, u.LanguageCode,
		u.EscapedName(), h.svc.CFG.Name, h.svc.CFG.NewsChannel, u.Balance, u.BalanceUSD.StringFixedBank(5))
	if u.PremiumExpiresAt != nil && u.PremiumExpiresAt.After(time.Now()) {
		txt += "\n" + langf(premiumText, u.LanguageCode, u.PremiumExpiresAt.Format("2006.01.02 15:04:05"))
	}
	if ctx.Update.CallbackQuery != nil {
		return ctx.EditMessageTextVoid(txt,
			&illuminate.EditMessageTextOpts{
				ReplyMarkup: menu.InlineKeyboardMarkup,
			},
		)
	}

	return ctx.ReplyWithMenuVoid(txt, menu)
}

func (h *Handler) Exit(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if err := h.svc.SetUserBotState(u, ""); err != nil {
		return ctx.ReplyVoid("Unexpected error, try again later")
	}
	return h.Start(ctx)
}

func (h *Handler) MyChatMember(ctx *router.Context) error {
	newChatMemberStatus := ctx.Update.MyChatMember.NewChatMember.GetStatus()
	if ctx.Update.MyChatMember.Chat.Type == illuminate.ChatTypeChannel {
		_, err := h.svc.CreateOrUpdateChannel(&ctx.Update.MyChatMember.Chat, newChatMemberStatus != "administrator")
		if err != nil {
			return err
		}
		return nil
	}

	var isStopped bool
	if newChatMemberStatus == "left" || newChatMemberStatus == "kicked" {
		isStopped = true
	} else {
		isStopped = false
	}
	u := getUserFromContext(ctx.Context)
	return h.svc.UpdateUserStoppedStatus(u, isStopped)
}

func (h *Handler) ChatMember(ctx *router.Context) error {
	newChatMemberStatus := ctx.Update.ChatMember.NewChatMember.GetStatus()
	newChatMemberUser := ctx.Update.ChatMember.NewChatMember.GetUser()
	reward, err := h.svc.ProcessChannelChatMember(
		ctx.Update.ChatMember.Chat.Id,
		newChatMemberUser.Id,
		newChatMemberStatus,
		ctx.Update.ChatMember.InviteLink,
	)
	if err != nil {
		return err
	}
	if reward == nil {
		return nil
	}

	if reward.IsReward {
		_, _ = h.bot.SendMessage(
			newChatMemberUser.Id,
			fmt.Sprintf("💠 Sponsor (%s) reward earned. You received <b>🪙 %d.</b>",
				reward.Channel.Title, reward.Channel.Reward,
			),
			&illuminate.SendMessageOpts{
				ParseMode: illuminate.ParseModeHTML,
			})
	} else {
		_, _ = h.bot.SendMessage(
			newChatMemberUser.Id,
			fmt.Sprintf("💠 Sponsor (%s) fine received. Your fine <b>🪙 -%d.</b>",
				reward.Channel.Title, reward.Channel.Reward*2,
			),
			&illuminate.SendMessageOpts{
				ParseMode: illuminate.ParseModeHTML,
			})
	}

	return nil
}

func (h *Handler) OnMessage(ctx *router.Context) error {
	return ctx.ReplyVoid("undefined command")
}
