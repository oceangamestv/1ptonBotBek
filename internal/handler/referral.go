package handler

import (
	"fmt"
	"github.com/kbgod/coinbot/internal/service"
	"github.com/kbgod/illuminate/router"
)

var referralText = map[string]string{
	"en": "🔗 Referral" +
		"\n\nInvite your friends and get rewards." +
		"\n\nYour referral link:" +
		"\n<code>%s</code>" +
		"\n\nFor each friend you invite, you will receive <code>$%s</code>." +
		"\nFor each friend you invite, you will receive 50%% of their taps (from 2 level)." +
		"\n\nReferrals count: <b>%d</b>" +
		"\nEarned: <b>%d</b> 🪙" +
		"\nEarned USD: <b>$%s</b>",
	"ru": "🔗 Реферальная программа" +
		"\n\nПриглашайте друзей и получайте вознаграждение." +
		"\n\nВаша реферальная ссылка:" +
		"\n<code>%s</code>" +
		"\n\nЗа каждого друга, которого вы пригласите, вы получите <code>$%s</code>." +
		"\nЗа каждого друга, которого вы пригласите, вы получите 50%% от их кликов (со 2 уровня)." +
		"\n\nКоличество рефералов: <b>%d</b>" +
		"\nЗаработано: <b>%d</b> 🪙" +
		"\nЗаработано USD: <b>$%s</b>",
	"uk": "🔗 Реферальна програма" +
		"\n\nЗапрошуйте друзів та отримуйте винагороду." +
		"\n\nВаше реферальне посилання:" +
		"\n<code>%s</code>" +
		"\n\nЗа кожного друга, якого ви запросите, ви отримаєте <code>$%s</code>." +
		"\nЗа кожного друга, якого ви запросите, ви отримаєте 50%% від їх кліків (з 2 рівня)." +
		"\n\nКількість рефералів: <b>%d</b>" +
		"\nЗароблено: <b>%d</b> 🪙" +
		"\nЗароблено USD: <b>$%s</b>",
}

func (h *Handler) Referral(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)

	link := fmt.Sprintf("https://t.me/%s?start=r_%d", h.bot.Username, u.ID)
	//txt := "<b>🔗 Referral</b>\n\n"
	//txt += "Invite your friends and get rewards.\n\n"
	//txt += "Your referral link:\n"
	//txt += fmt.Sprintf("<code>%s</code>\n\n", link)
	//txt += fmt.Sprintf("For each friend you invite, you will receive <code>$%s</code>.\n", service.ReferralRegisterReward.StringFixedBank(3))
	//txt += "For each friend you invite, you will receive 50% of their taps (from 2 level).\n\n"
	//
	refCount, err := h.svc.GetReferralsCount(u)
	if err != nil {
		return err
	}
	//txt += fmt.Sprintf("Referrals count: <b>%d</b>\n", refCount)
	//txt += fmt.Sprintf("Earned: <b>%d</b> 🪙\n", u.ReferralProfit)
	//txt += fmt.Sprintf("Earned USD: <b>$%s</b>", u.ReferralProfitUSD.StringFixedBank(3))

	return ctx.ReplyVoid(
		langf(
			referralText, u.LanguageCode,
			link,
			service.ReferralRegisterReward.StringFixedBank(3),
			refCount,
			u.ReferralProfit,
			u.ReferralProfitUSD.StringFixedBank(3)))
}
