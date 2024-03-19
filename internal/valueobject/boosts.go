package valueobject

import "time"

type DailyBooster struct {
	Coin   int64     `json:"coin"`
	Energy int32     `json:"energy"`
	NextAt time.Time `json:"next_at"`
}
