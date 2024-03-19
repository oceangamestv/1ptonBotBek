package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/valueobject"
	"github.com/kbgod/coinbot/pkg/randomizer"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var (
	ErrInsufficientBalance    = errors.New("insufficient balance")
	ErrInsufficientBalanceUSD = errors.New("insufficient balance in usd")
	ErrEnergySpeedLimit       = errors.New("energy speed limit reached")
)

func (s *Service) BuyMultitap(u *entity.User) error {
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return err
	}
	if u.Balance < u.MineLevelPrice() {
		return ErrInsufficientBalance
	}
	if err := s.db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"balance":    gorm.Expr("balance - ?", u.MineLevelPrice()),
			"mine_level": gorm.Expr("mine_level + 1"),
		}).Error; err != nil {
		return err
	}
	u.Balance -= u.MineLevelPrice()
	u.MineLevel++
	return nil
}

func (s *Service) BuyRechargeSpeed(u *entity.User) error {
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return err
	}
	if u.Balance < u.EnergyLevelPrice() {
		return ErrInsufficientBalance
	}
	if u.EnergyLevel >= 4 {
		return ErrEnergySpeedLimit
	}
	if err := s.db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"balance":      gorm.Expr("balance - ?", u.EnergyLevelPrice()),
			"energy_level": gorm.Expr("energy_level + 1"),
		}).Error; err != nil {
		return err
	}
	u.Balance -= u.EnergyLevelPrice()
	u.EnergyLevel++
	return nil
}

func (s *Service) BuyMaxEnergyLimit(u *entity.User) error {
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return err
	}
	if u.Balance < u.MaxEnergyPrice() {
		return ErrInsufficientBalance
	}
	if err := s.db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"balance":          gorm.Expr("balance - ?", u.MaxEnergyPrice()),
			"max_energy_level": gorm.Expr("max_energy_level + 1"),
		}).Error; err != nil {
		return err
	}
	u.Balance -= u.MaxEnergyPrice()
	u.MaxEnergyLevel++
	return nil
}

var AutoFarmerPrice = int64(100000)

func (s *Service) BuyAutoFarmer(u *entity.User) error {
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return err
	}
	if u.Balance < AutoFarmerPrice {
		return ErrInsufficientBalance
	}
	if err := s.db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"balance":     gorm.Expr("balance - ?", AutoFarmerPrice),
			"auto_farmer": true,
		}).Error; err != nil {
		return err
	}
	u.Balance -= AutoFarmerPrice
	u.AutoFarmer = true
	return nil
}

func (s *Service) OpenDailyBooster(u *entity.User) (*valueobject.DailyBooster, error) {
	if u.DailyBoosterAvailableAt.After(time.Now()) {
		return nil, errors.New("booster is not available yet")
	}
	if err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).Take(u).Error; err != nil {
		return nil, err
	}
	// crypto safe random number generator
	bigCoin, err := randomizer.GenerateRandomNumber(1, 100)
	if err != nil {
		return nil, err
	}
	var maxCoin int64
	if bigCoin >= 70 {
		maxCoin = 5000
	} else {
		maxCoin = 1200
	}
	coin, err := randomizer.GenerateRandomNumber(500, maxCoin)
	if err != nil {
		return nil, err
	}

	energy, err := randomizer.GenerateRandomNumber(500, int64(u.MaxEnergy()))
	if err != nil {
		return nil, err
	}
	if u.CurrentEnergy()+int32(energy) > u.MaxEnergy() {
		energy = int64(u.MaxEnergy() - u.CurrentEnergy())
	}
	fmt.Println(time.Now().UTC())
	fmt.Println(time.Now().Add(24 * time.Hour).UTC())
	fmt.Println(now.With(time.Now().Add(24 * time.Hour).UTC()).BeginningOfDay())
	fmt.Println(now.With(time.Now().Add(24 * time.Hour).UTC()).BeginningOfDay().UTC())

	next := now.With(time.Now().UTC().Add(24 * time.Hour)).BeginningOfDay().UTC()
	if err := s.db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"balance":                    gorm.Expr("balance + ?", coin),
			"energy":                     u.CurrentEnergy() + int32(energy),
			"daily_booster_available_at": next,
		}).Error; err != nil {
		return nil, err
	}

	return &valueobject.DailyBooster{
		Coin:   coin,
		Energy: int32(energy),
		NextAt: next,
	}, nil
}
