package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/kbgod/coinbot/internal/entity"
	"gorm.io/gorm"
	"time"
)

func initMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "init",
		Migrate: func(tx *gorm.DB) error {
			type Promo struct {
				ID        int64  `gorm:"primaryKey"`
				Name      string `gorm:"uniqueIndex"`
				Price     int64
				Charge    int64
				CreatedAt time.Time
			}
			type User struct {
				ID                         int64 `gorm:"primaryKey;autoIncrement:false"`
				FirstName                  string
				Username                   string
				Role                       entity.UserRole
				BanReason                  *string
				PromoID                    *int64
				BotState                   *string
				BotStateContext            *string
				Balance                    int64 `gorm:"default:0"`
				Energy                     int32 `gorm:"default:1000"`
				EnergyLevel                int32 `gorm:"default:1"` // speed
				MaxEnergyLevel             int32 `gorm:"default:0"`
				MineLevel                  int32 `gorm:"default:1"`
				AvatarURL                  *string
				WebAppAccessToken          *string
				RefererID                  *int64
				ReferralProfit             int64 `gorm:"default:0"`
				WebAppAccessTokenExpiresAt *time.Time
				StoppedAt                  *time.Time
				BannedAt                   *time.Time
				LastEnergyAt               time.Time `gorm:"default:now()"`
				LastMineAt                 time.Time `gorm:"default:now()"`
				CreatedAt                  time.Time
			}
			type Channel struct {
				ID         int64 `gorm:"primaryKey;autoIncrement:false"`
				Title      string
				InviteLink string
				Balance    int64 `gorm:"default:0"`
				Reward     int64 `gorm:"default:1000"`
				Activated  bool  `gorm:"default:false;index"`
				StoppedAt  *time.Time
			}

			type ChannelMember struct {
				ChannelID int64 `gorm:"primaryKey"`
				UserID    int64 `gorm:"primaryKey"`
				Status    string
				UpdatedAt time.Time
				CreatedAt time.Time
			}
			return tx.AutoMigrate(
				&Promo{},
				&User{},
				&Channel{},
				&ChannelMember{},
			)
		},
	}
}
