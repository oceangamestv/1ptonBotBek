package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

func premiumMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "4_premium",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				BalanceUSD        decimal.Decimal `gorm:"type:decimal(20,5);default:0"`
				ReferralProfitUSD decimal.Decimal `gorm:"type:decimal(20,5);default:0"`
				PremiumExpiresAt  *time.Time
				PhoneNumber       *string
			}
			type Promo struct {
				PremiumMinutes int32 `gorm:"default:0"`
			}
			return tx.AutoMigrate(&User{}, &Promo{})
		},
		Rollback: func(tx *gorm.DB) error {
			type User struct{}
			type Promo struct{}
			if err := tx.Migrator().DropColumn(&User{}, "balance_usd"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&User{}, "referral_profit_usd"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&User{}, "premium_expires_at"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&User{}, "phone_number"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Promo{}, "premium_minutes"); err != nil {
				return err
			}

			return nil
		},
	}
}
