package api

import "github.com/gofiber/fiber/v2"

func (h *handler) leaderboard(ctx *fiber.Ctx) error {
	leaderboard, err := h.svc.Leaderboard(getUserFromCtx(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(leaderboard)
}

func (h *handler) dailyLeaderboard(ctx *fiber.Ctx) error {
	leaderboard, err := h.svc.DailyLeaderboard(getUserFromCtx(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(leaderboard)
}

func (h *handler) monthlyLeaderboard(ctx *fiber.Ctx) error {
	leaderboard, err := h.svc.MonthlyLeaderboard(getUserFromCtx(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(leaderboard)
}
