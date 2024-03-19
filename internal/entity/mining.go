package entity

import "time"

type MiningPerDay struct {
	Date   time.Time `gorm:"primaryKey"`
	UserID int64     `gorm:"primaryKey"`
	User   *User
	League uint8
	Mined  int64
}

type MiningPerMonth struct {
	Date   time.Time `gorm:"primaryKey"`
	UserID int64     `gorm:"primaryKey"`
	User   *User
	League uint8
	Mined  int64
}
