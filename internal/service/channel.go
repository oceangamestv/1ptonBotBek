package service

import (
	"errors"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/illuminate"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

func (s *Service) CreateOrUpdateChannel(chat *illuminate.Chat, stopped bool) (*entity.Channel, error) {
	channelEntity := entity.Channel{
		ID:    chat.ID,
		Title: chat.Title,
	}
	if stopped {
		n := time.Now()
		channelEntity.StoppedAt = &n
	}

	if err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "stopped_at"}),
	}).Create(&channelEntity).Error; err != nil {
		return nil, err
	}

	return &channelEntity, nil
}

func (s *Service) GetActiveChannels(userID int64) ([]entity.Channel, error) {
	var channels []entity.Channel
	if err := s.db.
		Where("activated = true and stopped_at IS NULL").
		Where("balance > reward").
		Order("reward desc").
		Find(&channels).Error; err != nil {
		return nil, err
	}

	return channels, nil
}

type ChannelReward struct {
	Channel  *entity.Channel
	IsReward bool
}

func (s *Service) ProcessChannelChatMember(
	chatID int64, userID int64, status string, inviteLink *illuminate.ChatInviteLink,
) (*ChannelReward, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	channel := &entity.Channel{}
	err := tx.Take(channel, "id = ?", chatID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if channel.StoppedAt != nil || channel.Balance <= channel.Reward {
		return nil, nil
	}
	if inviteLink != nil && status == "member" && !strings.HasPrefix(
		channel.InviteLink, strings.TrimSuffix(inviteLink.InviteLink, "...")) {
		return nil, nil
	}

	chatMember := &entity.ChannelMember{}
	err = tx.Take(chatMember, "channel_id = ? AND user_id = ?", chatID, userID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var rewardedChannel *ChannelReward
	if errors.Is(err, gorm.ErrRecordNotFound) { // зняти баланс з каналу і нарахувати користувачу бонус
		if status == "member" {
			if err := tx.
				Model(&entity.Channel{}).
				Where("id = ?", chatID).
				Update("balance", gorm.Expr("balance - reward")).Error; err != nil {
				return nil, err
			}
			if err := tx.
				Model(&entity.User{}).
				Where("id = ?", userID).
				Update("balance", gorm.Expr("balance + ?", channel.Reward)).Error; err != nil {
				return nil, err
			}

			chatMember = &entity.ChannelMember{
				ChannelID: chatID,
				UserID:    userID,
				Status:    status,
			}
			if err := tx.Create(chatMember).Error; err != nil {
				return nil, err
			}
			rewardedChannel = &ChannelReward{
				Channel:  channel,
				IsReward: true,
			}
		}
	} else { // користувач вже був в базі, і якщо користувач був мембером, і відписався, то зняти бонус
		if chatMember.Status == "member" && (status == "left" || status == "kicked") {
			if err := tx.
				Model(&entity.User{}).
				Where("id = ?", userID).
				Update("balance", gorm.Expr("balance - ?", channel.Reward*2)).Error; err != nil {
				return nil, err
			}
			if err := tx.
				Model(&entity.ChannelMember{}).
				Where("channel_id = ? AND user_id = ?", chatID, userID).
				Update("status", status).Error; err != nil {
				return nil, err
			}

			rewardedChannel = &ChannelReward{
				Channel:  channel,
				IsReward: false,
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return rewardedChannel, nil
}
