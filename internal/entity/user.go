package entity

import (
	"github.com/kbgod/coinbot/pkg/tgutil"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type League struct {
	MustReachBalance int64
	Name             string
}

var Leagues = []League{
	{MustReachBalance: 1000, Name: "🪵 Wood"},
	{MustReachBalance: 5000, Name: "🥉 Bronze"},
	{MustReachBalance: 100_000, Name: "🥈 Silver"},
	{MustReachBalance: 1_000_000, Name: "🥇 Gold"},
	{MustReachBalance: 10_000_000, Name: "💎 Diamond"},
	{MustReachBalance: 100_000_000, Name: "🌟 Platinum"},
}

type User struct {
	ID                         int64 `gorm:"primaryKey;autoIncrement:false"`
	FirstName                  string
	Username                   string
	Role                       UserRole
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
	ReferralProfit             int64           `gorm:"default:0"`
	ReferralProfitUSD          decimal.Decimal `gorm:"type:decimal(20,5);default:0"`
	LanguageCode               string
	AutoFarmer                 bool
	AutoFarmerProfit           int64           `gorm:"-"`
	BalanceUSD                 decimal.Decimal `gorm:"type:decimal(20,5);default:0"`
	PhoneNumber                *string
	PremiumExpiresAt           *time.Time
	WebAppAccessTokenExpiresAt *time.Time
	StoppedAt                  *time.Time
	BannedAt                   *time.Time
	LastEnergyAt               time.Time `gorm:"default:now()"`
	LastMineAt                 time.Time `gorm:"default:now()"`
	DailyBoosterAvailableAt    time.Time `gorm:"default:now()"`
	League                     uint8
	CreatedAt                  time.Time
}

func (u *User) IsPremium() bool {
	return u.PremiumExpiresAt != nil && u.PremiumExpiresAt.After(time.Now())
}

func (u *User) Identity() string {
	if u.Username != "" {
		return u.Username
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	return "unknown"
}

func (u *User) MineLevelPrice() int64 {
	// 1000 coins per level
	return int64(math.Pow(2, float64(u.MineLevel-1))) * 1000
}

func (u *User) EnergyLevelPrice() int64 {
	return int64(math.Pow(2, float64(u.EnergyLevel-1))) * 1000
}

func (u *User) MaxEnergyPrice() int64 {
	return int64(math.Pow(2, float64(u.MaxEnergyLevel))) * 1000
}

func (u *User) MaxEnergy() int32 {
	// +500 energy per level
	return 1000 + u.MaxEnergyLevel*500
}

func (u *User) CurrentEnergy() int32 {
	energy := int32(time.Since(u.LastMineAt).Seconds())
	energy *= u.EnergyLevel
	if energy+u.Energy > u.MaxEnergy() {
		energy = u.MaxEnergy() - u.Energy
	}
	return u.Energy + energy
}

func (u *User) EscapedName() string {
	return tgutil.Escape(u.FirstName)
}
