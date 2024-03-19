package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/service"
)

type boostsResponse struct {
	CurrentMineLevel   int32 `json:"current_mine_level"`
	MineLevelPrice     int64 `json:"mine_level_price"`
	CurrentEnergyLevel int32 `json:"current_energy_level"`
	EnergyLevelPrice   int64 `json:"energy_level_price"`
	CurrentMaxEnergy   int32 `json:"current_max_energy"`
	MaxEnergyPrice     int64 `json:"max_energy_price"`
	AutoFarmerPrice    int64 `json:"auto_farmer_price"`
}

func newBoostsResponse(user *entity.User) *boostsResponse {
	return &boostsResponse{
		CurrentMineLevel:   user.MineLevel,
		MineLevelPrice:     user.MineLevelPrice(),
		CurrentEnergyLevel: user.EnergyLevel,
		EnergyLevelPrice:   user.EnergyLevelPrice(),
		CurrentMaxEnergy:   user.MaxEnergy(),
		MaxEnergyPrice:     user.MaxEnergyPrice(),
		AutoFarmerPrice:    service.AutoFarmerPrice,
	}
}

func (h *handler) boosts(ctx *fiber.Ctx) error {
	user := getUserFromCtx(ctx)

	return ctx.JSON(newBoostsResponse(user))
}

type userWithBoostsResponse struct {
	User   *userResponse   `json:"user"`
	Boosts *boostsResponse `json:"boosts"`
}

func newUserWithBoostsResponse(user *entity.User) *userWithBoostsResponse {
	return &userWithBoostsResponse{
		User:   newUserResponse(user),
		Boosts: newBoostsResponse(user),
	}
}

type purchaseBoostRequest struct {
	Boost string `json:"boost"`
}

func (h *handler) purchaseBoost(ctx *fiber.Ctx) error {
	req := new(purchaseBoostRequest)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	user := getUserFromCtx(ctx)
	var err error
	switch req.Boost {
	case "multitap":
		err = h.svc.BuyMultitap(user)
	case "energy":
		err = h.svc.BuyRechargeSpeed(user)
	case "max_energy":
		err = h.svc.BuyMaxEnergyLimit(user)
	case "auto_farmer":
		err = h.svc.BuyAutoFarmer(user)
	}
	if err != nil && errors.Is(err, service.ErrInsufficientBalance) {
		return newMessageResponse(ctx, fiber.StatusForbidden, "insufficient balance")
	} else if err != nil {
		return err
	}

	return ctx.JSON(newUserWithBoostsResponse(user))
}

func (h *handler) openDailyBooster(ctx *fiber.Ctx) error {
	user := getUserFromCtx(ctx)
	booster, err := h.svc.OpenDailyBooster(user)
	if err != nil {
		return err
	}

	return ctx.JSON(booster)
}
