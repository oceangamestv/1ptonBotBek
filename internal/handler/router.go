package handler

import (
	"github.com/kbgod/illuminate/router"
	"strings"
)

func (h *Handler) initRoutes() *router.Router {
	botRouter := router.New(h.bot)
	botRouter.Use(h.Recovery)
	botRouter.Use(h.ErrorHandler)
	botRouter.Use(h.CallbackQueryAutoAnswer)
	botRouter.Use(h.UserMiddleware)

	// global events (accessible from any state)
	botRouter.On(onMyChatMember, h.MyChatMember)
	botRouter.On(onChatMember, h.ChatMember)
	botRouter.OnStart(h.Start)
	botRouter.OnCommand("set_my_name", h.SetUserName)
	botRouter.On(onWebAppData("earn"), h.EarnPage)
	botRouter.OnCommand("boosts", h.Boosts)
	botRouter.OnCommand("referral", h.Referral)
	botRouter.OnCommand("premium", h.Premium)
	botRouter.OnCallbackPrefix("premium", h.Premium)
	botRouter.OnCallbackPrefix("buy_premium/", h.BuyPremium)
	botRouter.OnCommand("multitap", h.MultitapPage)
	botRouter.OnCallbackPrefix("boosts", h.Boosts)
	botRouter.OnCallbackPrefix("referral", h.Referral)
	botRouter.OnCallbackPrefix("multitap_page", h.MultitapPage)
	botRouter.OnCallbackPrefix("buy_multitap", h.BuyMultitap)
	botRouter.OnCallbackPrefix("recharge_speed", h.RechargingSpeedPage)
	botRouter.OnCallbackPrefix("buy_recharge_speed", h.BuyRechargingSpeed)
	botRouter.OnCallbackPrefix("energy_limit", h.EnergyLimitPage)
	botRouter.OnCallbackPrefix("buy_energy_limit", h.BuyMaxEnergy)
	botRouter.OnCallbackPrefix("admin_panel", h.EnterAdminPanel)

	adminRouter := botRouter.UseState("admin_panel")
	adminRouter.OnCommand("exit", h.Exit)
	adminRouter.OnMessage(h.AdminMainMenu)
	adminRouter.OnCallbackPrefix("sending", h.EnterSendingPanel)

	sendingRouter := botRouter.UseState("sending")
	sendingRouter.On(onKeyboardFormat, h.WithKeyboard)
	sendingRouter.OnCommand("exit", h.ExitFromSending)
	sendingRouter.OnCommand("skip", h.WithoutKeyboard)
	sendingRouter.OnMessage(h.PrepareMessage)
	sendingRouter.OnCallbackPrefix("send", h.Send)

	h.StartSenderWorker()

	// this handler will be called only if routes with states and type OnMessage are not defined
	botRouter.OnMessage(h.Start)

	return botRouter
}

func onKeyboardFormat(ctx *router.Context) bool {
	if ctx.Update.Message == nil {
		return false
	}
	if ctx.Update.Message.Text == "" {
		return false
	}

	lines := strings.Split(ctx.Update.Message.Text, "\n")
	for _, line := range lines {
		if !strings.Contains(line, " - ") {
			return false
		}
	}

	return true
}

func onWebAppData(data string) func(ctx *router.Context) bool {
	return func(ctx *router.Context) bool {
		if ctx.Update.Message == nil {
			return false
		}
		if ctx.Update.Message.WebAppData == nil {
			return false
		}
		return ctx.Update.Message.WebAppData.Data == data
	}
}

func onMyChatMember(ctx *router.Context) bool {
	return ctx.Update.MyChatMember != nil
}

func onChatMember(ctx *router.Context) bool {
	return ctx.Update.ChatMember != nil
}
