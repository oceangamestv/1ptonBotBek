package entity

import "time"

type Channel struct {
	ID           int64 `gorm:"primaryKey;autoIncrement:false"`
	Title        string
	InviteLink   string
	Balance      int64 `gorm:"default:0"`
	Reward       int64 `gorm:"default:1000"`
	Activated    bool  `gorm:"default:false;index"`
	LanguageCode string
	StoppedAt    *time.Time
}

type ChannelMember struct {
	ChannelID int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"primaryKey"`
	Status    string
	UpdatedAt time.Time
	CreatedAt time.Time
}
