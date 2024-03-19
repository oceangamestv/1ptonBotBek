package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/service"
	"gorm.io/gorm"
)

type handler struct {
	svc *service.Service
}

func newHandler(svc *service.Service) *handler {
	return &handler{svc: svc}
}

const userKey = "user"

func (h *handler) userMiddleware(ctx *fiber.Ctx) error {
	accessToken := ctx.Get("x-api-key")
	if accessToken == "" {
		return newMessageResponse(ctx, fiber.StatusUnauthorized, "missing x-api-key header")
	}

	user, err := h.svc.GetUserByAccessToken(accessToken)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return newMessageResponse(ctx, fiber.StatusUnauthorized, "invalid x-api-key header")
	}
	if err != nil {
		return err
	}

	ctx.Locals(userKey, user)

	return ctx.Next()
}

func getUserFromCtx(ctx *fiber.Ctx) *entity.User {
	return ctx.Locals(userKey).(*entity.User)
}
