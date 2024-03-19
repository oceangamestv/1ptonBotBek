package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"time"
)

func dailyBoostersMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "daily boosters",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				DailyBoosterAvailableAt time.Time `gorm:"default:now()"`
			}
			type DailyBoosterHistory struct {
				UserID int64     `gorm:"primaryKey"`
				Date   time.Time `gorm:"primaryKey;type:date"`
			}
			return tx.AutoMigrate(&User{})
		},
	}
}
