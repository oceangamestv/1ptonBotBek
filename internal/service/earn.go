package service

import "github.com/shopspring/decimal"

type PremiumPack struct {
	Days  int
	Price decimal.Decimal
}

var PremiumPacks = []PremiumPack{
	{Days: 1, Price: decimal.NewFromInt(1)},
	{Days: 3, Price: decimal.NewFromFloat(2.5)},
	{Days: 7, Price: decimal.NewFromInt(5)},
	{Days: 30, Price: decimal.NewFromInt(20)},
	{Days: 90, Price: decimal.NewFromInt(50)},
	{Days: 180, Price: decimal.NewFromInt(100)},
	{Days: 365, Price: decimal.NewFromInt(180)},
}

var (
	ReferralRegisterReward = decimal.NewFromFloat(0.03)
)
