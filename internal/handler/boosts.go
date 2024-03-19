package handler

import (
	"errors"
	"fmt"
	"github.com/kbgod/coinbot/internal/service"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
)

func (h *Handler) Boosts(ctx *router.Context) error {
	txt := "<b>🚀 Boosts</b>\n\n"
	txt += "Boosts are special items that can help you to mine more efficiently.\n\n"

	menu := illuminate.NewInlineMenu()
	menu.Row().CallbackBtn("👆 Multitap", "multitap_page")
	menu.Row().CallbackBtn("⚡️ Recharging speed", "recharge_speed")
	menu.Row().CallbackBtn("🔋 Energy limit", "energy_limit")
	if ctx.Update.CallbackQuery != nil {
		return ctx.EditMessageTextVoid(txt, &illuminate.EditMessageTextOpts{
			ReplyMarkup: menu.InlineKeyboardMarkup,
		})
	}
	return ctx.ReplyWithMenuVoid(txt, menu)
}

func (h *Handler) MultitapPage(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	txt := "<b>👆 Multitap</b>\n\n"
	txt += "Multitap increases amount of coins you can earn per one tap.\n\n"
	txt += fmt.Sprintf("Current level: <code>%d</code>\n", u.MineLevel)
	txt += fmt.Sprintf("Your balance: <code>%d 🪙</code>\n\n", u.Balance)
	txt += "<code>+1 per tap for each level.</code>"

	menu := illuminate.NewInlineMenu()

	menu.
		Row().
		CallbackBtn(
			fmt.Sprintf("Get for %d 🪙", u.MineLevelPrice()),
			"buy_multitap",
		)
	menu.Row().CallbackBtn("« Back", "boosts")
	return ctx.EditMessageTextVoid(txt, &illuminate.EditMessageTextOpts{
		ReplyMarkup: menu.InlineKeyboardMarkup,
	})
}

func (h *Handler) BuyMultitap(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.Balance < u.MineLevelPrice() {
		diff := u.MineLevelPrice() - u.Balance
		return ctx.AnswerAlertVoid(
			fmt.Sprintf("You need %d more 🪙 to buy this multitap level", diff),
		)
	}
	err := h.svc.BuyMultitap(u)
	if err != nil && errors.Is(err, service.ErrInsufficientBalance) {
		diff := u.MineLevelPrice() - u.Balance
		return ctx.AnswerAlertVoid(
			fmt.Sprintf("❌ You need %d more 🪙 to buy this multitap level", diff),
		)
	} else if err != nil {
		return err
	}
	if err := h.MultitapPage(ctx); err != nil {
		return err
	}
	return ctx.AnswerAlertVoid(fmt.Sprintf("✅ You bought %d multitap level.", u.MineLevel))
}

func (h *Handler) RechargingSpeedPage(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	txt := "<b>⚡️ Recharging speed</b>\n\n"
	txt += "Recharging speed increases amount of energy points (EP) you receive every 2 seconds\n\n"
	txt += fmt.Sprintf("Current level: <code>%d EP per 2 sec.</code>\n", u.EnergyLevel)
	txt += fmt.Sprintf("Your balance: <code>%d 🪙</code>\n\n", u.Balance)
	txt += "<code>+1 EP per 2 seconds.</code>"

	menu := illuminate.NewInlineMenu()

	menu.
		Row().
		CallbackBtn(
			fmt.Sprintf("Get for %d 🪙", u.EnergyLevelPrice()),
			"buy_recharge_speed",
		)
	menu.Row().CallbackBtn("« Back", "boosts")
	return ctx.EditMessageTextVoid(txt, &illuminate.EditMessageTextOpts{
		ReplyMarkup: menu.InlineKeyboardMarkup,
	})
}

func (h *Handler) BuyRechargingSpeed(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.Balance < u.EnergyLevelPrice() {
		diff := u.EnergyLevelPrice() - u.Balance
		return ctx.AnswerAlertVoid(
			fmt.Sprintf("You need %d more 🪙 to buy this recharging speed level", diff),
		)
	}
	err := h.svc.BuyRechargeSpeed(u)
	if err != nil && errors.Is(err, service.ErrInsufficientBalance) {
		diff := u.EnergyLevelPrice() - u.Balance
		return ctx.AnswerAlertVoid(
			fmt.Sprintf("❌ You need %d more 🪙 to buy this recharging speed level", diff),
		)
	} else if err != nil {
		return err
	}
	if err := h.RechargingSpeedPage(ctx); err != nil {
		return err
	}
	return ctx.AnswerAlertVoid(fmt.Sprintf("✅ You bought %d recharging speed level.", u.EnergyLevel))
}

func (h *Handler) EnergyLimitPage(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	txt := "<b>🔋 Energy limit</b>\n\n"
	txt += "Energy limit increases your maximum energy points (EP)\n\n"
	txt += fmt.Sprintf("Current capacity: <code>%d EP</code>\n", u.MaxEnergy())
	txt += fmt.Sprintf("Your balance: <code>%d 🪙</code>\n\n", u.Balance)
	txt += "<code>+500 EP per level.</code>"

	menu := illuminate.NewInlineMenu()

	menu.
		Row().
		CallbackBtn(
			fmt.Sprintf("Get +500 for %d 🪙", u.MaxEnergyPrice()),
			"buy_energy_limit",
		)
	menu.Row().CallbackBtn("« Back", "boosts")
	return ctx.EditMessageTextVoid(txt, &illuminate.EditMessageTextOpts{
		ReplyMarkup: menu.InlineKeyboardMarkup,
	})
}

func (h *Handler) BuyMaxEnergy(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	if u.Balance < u.MaxEnergyPrice() {
		diff := u.MaxEnergyPrice() - u.Balance
		return ctx.AnswerAlertVoid(
			fmt.Sprintf("You need %d more 🪙 to buy this energy limit level", diff),
		)
	}
	err := h.svc.BuyMaxEnergyLimit(u)
	if err != nil && errors.Is(err, service.ErrInsufficientBalance) {
		diff := u.MaxEnergyPrice() - u.Balance
		return ctx.AnswerAlertVoid(
			fmt.Sprintf("❌ You need %d more 🪙 to buy this energy limit level", diff),
		)
	} else if err != nil {
		return err
	}
	if err := h.EnergyLimitPage(ctx); err != nil {
		return err
	}
	return ctx.AnswerAlertVoid(fmt.Sprintf("✅ You bought %d energy limit level.", u.MaxEnergyLevel))
}
