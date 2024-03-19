package handler

import (
	"fmt"
	"github.com/kbgod/illuminate"
)

func (h *Handler) initCommands() error {
	ok, err := h.bot.SetMyCommands([]illuminate.BotCommand{
		{
			Command:     "start",
			Description: "main menu",
		},
		{
			Command:     "premium",
			Description: "buy premium",
		},
		{
			Command:     "referral",
			Description: "your referral link",
		},
	}, &illuminate.SetMyCommandsOpts{
		Scope:        illuminate.BotCommandScopeDefault{},
		LanguageCode: "en",
	})
	if err != nil {
		return fmt.Errorf("set en commands: %w", err)
	}
	if !ok {
		return fmt.Errorf("set en commands: not ok")
	}

	ok, err = h.bot.SetMyCommands([]illuminate.BotCommand{
		{
			Command:     "start",
			Description: "главное меню",
		},
		{
			Command:     "premium",
			Description: "купить премиум",
		},
		{
			Command:     "referral",
			Description: "ваша реферальная ссылка",
		},
	}, &illuminate.SetMyCommandsOpts{
		Scope:        illuminate.BotCommandScopeDefault{},
		LanguageCode: "ru",
	})
	if err != nil {
		return fmt.Errorf("set commands: %w", err)
	}
	if !ok {
		return fmt.Errorf("set commands: not ok")
	}

	ok, err = h.bot.SetMyCommands([]illuminate.BotCommand{
		{
			Command:     "start",
			Description: "головне меню",
		},
		{
			Command:     "premium",
			Description: "купити преміум",
		},
		{
			Command:     "referral",
			Description: "ваше реферальне посилання",
		},
	}, &illuminate.SetMyCommandsOpts{
		Scope:        illuminate.BotCommandScopeDefault{},
		LanguageCode: "uk",
	})
	if err != nil {
		return fmt.Errorf("set uk commands: %w", err)
	}
	if !ok {
		return fmt.Errorf("set uk commands: not ok")
	}

	return nil
}
