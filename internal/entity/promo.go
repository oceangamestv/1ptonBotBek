package entity

import (
	"gorm.io/gorm"
	"time"
)

type Promo struct {
	ID             int64  `gorm:"primaryKey"`
	Name           string `gorm:"uniqueIndex"`
	Price          int64
	Charge         int64
	PremiumMinutes int32 `gorm:"default:0"`
	CreatedAt      time.Time
}

func GetPromoByName(tx *gorm.DB, name string) (*Promo, error) {
	p := &Promo{}
	return p, tx.Take(p, "name = ?", name).Error
}
