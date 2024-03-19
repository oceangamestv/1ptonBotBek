package service

import (
	"errors"
	"github.com/jinzhu/now"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/valueobject"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var (
	ErrInsufficientEnergy = errors.New("insufficient energy")
)

var (
	ErrMiningTooFast = errors.New("mining too fast")
)

func (s *Service) MineMany(u *entity.User, count int32) (*valueobject.MineManyResult, error) {
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return nil, err
	}
	//maximum 20 click per second
	diff := time.Since(u.LastMineAt.Truncate(time.Millisecond))
	minDuration := time.Duration(count) * time.Second / 30
	if diff < minDuration {
		return nil, ErrMiningTooFast
	}

	energy := int32(diff.Seconds())
	energy *= u.EnergyLevel

	if energy+u.Energy > u.MaxEnergy() {
		energy = u.MaxEnergy() - u.Energy
	}
	// 26 > 12 * 3
	if u.Energy+energy < u.MineLevel*count && count > 1 {
		// decrease count
		count = (u.Energy + energy) / u.MineLevel
	}
	if count < 1 || u.Energy+energy < u.MineLevel*count {
		return nil, ErrInsufficientEnergy
	}

	mul := int32(1)
	if u.IsPremium() {
		mul = 2
	}

	mined := u.MineLevel * mul * count

	if err := s.db.Model(u).Updates(map[string]any{
		"energy":       gorm.Expr("energy - ?", u.MineLevel*count-energy),
		"balance":      gorm.Expr("balance + ?", mined),
		"last_mine_at": time.Now().UTC(),
	}).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&entity.MiningPerDay{}).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "date"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"mined": gorm.Expr("mining_per_days.mined + ?", mined),
		}),
	}).Create(&entity.MiningPerDay{
		Date:   now.With(time.Now().UTC()).BeginningOfDay(),
		UserID: u.ID,
		League: u.League,
		Mined:  int64(mined),
	}).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&entity.MiningPerMonth{}).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "date"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"mined": gorm.Expr("mining_per_months.mined + ?", mined),
		}),
	}).Create(&entity.MiningPerDay{
		Date:   now.With(time.Now().UTC()).BeginningOfMonth(),
		UserID: u.ID,
		League: u.League,
		Mined:  int64(mined),
	}).Error; err != nil {
		return nil, err
	}
	if u.RefererID != nil && u.MineLevel > 1 {
		if err := s.db.
			Model(&entity.User{}).
			Where("id = ?", *u.RefererID).
			Updates(map[string]any{
				"balance":         gorm.Expr("balance + ?", count/2),
				"referral_profit": gorm.Expr("referral_profit + ?", count/2),
			}).Error; err != nil {
			return nil, err
		}
	}

	return &valueobject.MineManyResult{
		Balance:   u.Balance + int64(mined),
		Mined:     u.MineLevel * mul * count,
		NewEnergy: u.Energy + energy - u.MineLevel*count,
	}, nil
}

func (s *Service) MakeEnergy(u *entity.User) (int32, error) {
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return 0, err
	}
	dur := time.Since(u.LastEnergyAt)
	// energy point per 2 seconds
	energy := int32(dur.Seconds() / 2)
	energy *= u.EnergyLevel
	if energy+u.Energy > u.MaxEnergy() {
		energy = u.MaxEnergy() - u.Energy
	}
	if energy == 0 {
		return u.Energy, nil
	}

	if err := s.db.Model(u).Updates(map[string]any{
		"energy":         gorm.Expr("energy + ?", energy),
		"last_energy_at": time.Now().UTC(),
	}).Error; err != nil {
		return 0, err
	}

	return u.Energy + energy, nil
}
