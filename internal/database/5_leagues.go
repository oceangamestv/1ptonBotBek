package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"time"
)

func leaguesMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "5_leagues",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				League uint8
			}

			type MiningPerDay struct {
				Date   time.Time `gorm:"primaryKey;index:idx_mining_per_day_date_mined"`
				UserID int64     `gorm:"primaryKey"`
				League uint8
				Mined  int64 `gorm:"index:idx_mining_per_day_date_mined"`
			}

			type MiningPerMonth struct {
				Date   time.Time `gorm:"primaryKey;index:idx_mining_per_month_date_mined"`
				UserID int64     `gorm:"primaryKey"`
				League uint8
				Mined  int64 `gorm:"index:idx_mining_per_month_date_mined"`
			}

			if err := tx.AutoMigrate(&User{}, &MiningPerDay{}, &MiningPerMonth{}); err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			type User struct{}
			type MiningPerDay struct{}
			type MiningPerMonth struct{}

			if err := tx.Migrator().DropColumn(&User{}, "league"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&MiningPerDay{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&MiningPerMonth{}); err != nil {
				return err
			}
			return nil
		},
	}
}
