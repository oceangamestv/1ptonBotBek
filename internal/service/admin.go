package service

import (
	"fmt"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/valueobject"
	"time"
)

func (s *Service) GetStatistics() (*valueobject.Statistics, error) {
	var statistics valueobject.Statistics
	if err := s.db.Model(&entity.User{}).Count(&statistics.UsersCount).Error; err != nil {
		return nil, fmt.Errorf("get users count: %w", err)
	}
	if err := s.db.Model(&entity.User{}).Where("stopped_at IS NULL").Count(&statistics.ActiveUsersCount).Error; err != nil {
		return nil, fmt.Errorf("get active users count: %w", err)
	}
	if err := s.db.Model(&entity.User{}).Where("created_at > ?", time.Now().Truncate(24*time.Hour)).Count(&statistics.CreatedToday).Error; err != nil {
		return nil, fmt.Errorf("get created today: %w", err)
	}
	if err := s.db.Model(&entity.User{}).Where("stopped_at > ?", time.Now().Truncate(24*time.Hour)).Count(&statistics.StoppedToday).Error; err != nil {
		return nil, fmt.Errorf("get stopped today: %w", err)
	}
	if err := s.db.Model(&entity.User{}).Where("referer_id IS NOT NULL").Count(&statistics.CreatedByRef).Error; err != nil {
		return nil, fmt.Errorf("get created by ref: %w", err)
	}
	return &statistics, nil
}
