package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kbgod/coinbot/internal/service"
)

type mineManyRequest struct {
	Count int32 `json:"count"`
}

func (h *handler) mine(ctx *fiber.Ctx) error {
	req := &mineManyRequest{}
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	user := getUserFromCtx(ctx)
	if req.Count < 1 {
		return newMessageResponse(ctx, fiber.StatusBadRequest, "invalid count")
	}

	result, err := h.svc.MineMany(user, req.Count)
	if err != nil && errors.Is(err, service.ErrInsufficientEnergy) {
		return newMessageResponse(ctx, fiber.StatusForbidden, err.Error())
	} else if err != nil && errors.Is(err, service.ErrMiningTooFast) {
		return newMessageResponse(ctx, fiber.StatusTooManyRequests, err.Error())
	} else if err != nil {
		return err
	}

	return ctx.JSON(result)
}
