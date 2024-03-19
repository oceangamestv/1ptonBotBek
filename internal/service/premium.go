package service

import (
	"database/sql"
	"github.com/kbgod/coinbot/internal/entity"
	"gorm.io/gorm"
	"time"
)

func (s *Service) BuyPremium(u *entity.User, pack PremiumPack) error {
	tx := s.db.Begin(&sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Take(u, u.ID).Error; err != nil {
		return err
	}

	if u.BalanceUSD.LessThan(pack.Price) {
		return ErrInsufficientBalanceUSD
	}

	var premiumExpiresAt time.Time
	if u.PremiumExpiresAt != nil && u.PremiumExpiresAt.After(time.Now()) {
		premiumExpiresAt = u.PremiumExpiresAt.AddDate(0, 0, pack.Days)
	} else {
		premiumExpiresAt = time.Now().AddDate(0, 0, pack.Days)
	}

	if err := tx.Model(u).Updates(map[string]interface{}{
		"balance_usd":        gorm.Expr("balance_usd - ?", pack.Price),
		"premium_expires_at": premiumExpiresAt,
	}).Error; err != nil {
		return err
	}

	u.BalanceUSD = u.BalanceUSD.Sub(pack.Price)
	u.PremiumExpiresAt = &premiumExpiresAt

	return tx.Commit().Error
}
