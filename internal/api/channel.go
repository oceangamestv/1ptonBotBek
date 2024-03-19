package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kbgod/coinbot/internal/entity"
)

type channelResponse struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	InviteLink string `json:"invite_link"`
	Reward     int64  `json:"reward"`
}

func newChannelResponse(channel *entity.Channel) *channelResponse {
	return &channelResponse{
		ID:         channel.ID,
		Title:      channel.Title,
		InviteLink: channel.InviteLink,
		Reward:     channel.Reward,
	}
}

func newChannelsResponse(channels []entity.Channel) []*channelResponse {
	var resp []*channelResponse
	for _, channel := range channels {
		resp = append(resp, newChannelResponse(&channel))
	}
	return resp
}

func (h *handler) channels(ctx *fiber.Ctx) error {
	channels, err := h.svc.GetActiveChannels(getUserFromCtx(ctx).ID)
	if err != nil {
		return err
	}

	return ctx.JSON(newChannelsResponse(channels))
}
