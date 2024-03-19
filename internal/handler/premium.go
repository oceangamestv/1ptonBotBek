package handler

import (
	"fmt"
	"github.com/kbgod/coinbot/internal/service"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
	"strconv"
	"strings"
)

var premiumMenuText = map[string]string{
	"en": "🏵 Premium\n\n<b>Benefits:</b>\n<i>• X2 mining</i>\n<i>• A coin skin in the shape of a coconut</i>\n<i>• Name in the leaderboard in gold color</i>\n<i>• No ads</i>\n<i>• Support developers</i>\n\n<b>Your balance:</b> <code>$%s</code>\nFor top up write message: @kbgod\n\n⚠️ <i>After pressing the premium button, it will be purchased immediately</i>",
	"ru": "🏵 Премиум\n\n<b>Преимущества:</b>\n<i>• X2 майнинг</i>\n<i>• Скин на монету в виде кокоса</i>\n<i>• Имя в топе золотым цветом</i>\n<i>• Без рекламы</i>\n<i>• Поддержка разработчиков</i>\n\n<b>Ваш баланс:</b> <code>$%s</code>\nДля пополнения пишите: @kbgod\n\n⚠️ <i>После нажатия кнопки премиум будет сразу куплено</i>",
	"uk": "🏵 Преміум\n\n<b>Переваги:</b>\n<i>• X2 майнінг</i>\n<i>• Скін на монету у вигляді кокосу</i>\n<i>• Ім'я в топі золотим кольором</i>\n<i>• Немає реклами</i>\n<i>• Підтримка розробників</i>\n\n<b>Ваш баланс:</b> <code>$%s</code>\nДля поповнення пишіть: @kbgod\n\n⚠️ <i>Після натискання кнопки преміум буде придбано одразу</i>",
}

func (h *Handler) Premium(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)

	buttons := make([]illuminate.InlineKeyboardButton, 0, len(service.PremiumPacks))
	for id, pack := range service.PremiumPacks {
		buttons = append(buttons, illuminate.CallbackBtn(
			fmt.Sprintf("%dD - $%s", pack.Days, pack.Price), fmt.Sprintf("buy_premium/%d", id)),
		)
	}
	m := illuminate.NewInlineMenu().Fill(2, buttons...)

	return ctx.ReplyWithMenuVoid(langf(premiumMenuText, u.LanguageCode, u.BalanceUSD.StringFixedBank(5)), m)
}

var notEnoughBalanceText = map[string]string{
	"en": "❌ Not enough balance",
	"ru": "❌ Недостаточно средств",
	"uk": "❌ Недостатньо коштів",
}

var premiumPurchasedText = map[string]string{
	"en": "🏵 Premium purchased for %d days. Expires at %s",
	"ru": "🏵 Премиум куплен на %d дней. Истекает %s",
	"uk": "🏵 Преміум куплено на %d днів. Закінчується %s",
}

func (h *Handler) BuyPremium(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	id, err := strconv.Atoi(strings.TrimPrefix(ctx.Update.CallbackQuery.Data, "buy_premium/"))
	if err != nil {
		return err
	}
	if id > len(service.PremiumPacks) {
		return fmt.Errorf("premium pack not found: %d", id)
	}

	pack := service.PremiumPacks[id]

	if u.BalanceUSD.LessThan(pack.Price) {
		return ctx.AnswerAlertVoid(langf(notEnoughBalanceText, u.LanguageCode))
	}

	if err := h.svc.BuyPremium(u, pack); err != nil {
		return err
	}

	_ = ctx.AnswerAlertVoid(
		langf(
			premiumPurchasedText, u.LanguageCode,
			pack.Days, u.PremiumExpiresAt.UTC().Format("2006.01.02 15:04:05")+" UTC",
		),
	)

	return h.Start(ctx)
}
