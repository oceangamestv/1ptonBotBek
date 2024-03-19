package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func languagesMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3_languages",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				LanguageCode string
				AutoFarmer   bool
			}
			type Channel struct {
				LanguageCode string
			}
			return tx.AutoMigrate(&User{}, &Channel{})
		},
		Rollback: func(tx *gorm.DB) error {
			type User struct{}
			type Channel struct{}
			if err := tx.Migrator().DropColumn(&User{}, "language_code"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&User{}, "auto_farmer"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Channel{}, "language_code"); err != nil {
				return err
			}

			return nil
		},
	}
}
