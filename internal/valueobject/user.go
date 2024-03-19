package valueobject

import (
	"github.com/kbgod/coinbot/internal/entity"
	"time"
)

type PublicUser struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	Balance   int64   `json:"balance"`
	AvatarURL *string `json:"avatar_url"`
}

func NewPublicUserFromEntity(user entity.User) PublicUser {
	return PublicUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		Balance:   user.Balance,
		AvatarURL: user.AvatarURL,
	}
}

type UserScore struct {
	ID        int64   `json:"id"`
	Identity  string  `json:"identity"`
	Score     int64   `json:"score"`
	Position  int64   `json:"position"`
	AvatarURL *string `json:"avatar_url"`
	IsPremium bool    `json:"is_premium"`
}

func NewUserScoreFromEntity(user entity.User, position int64) UserScore {
	return UserScore{
		ID:        user.ID,
		Identity:  user.Identity(),
		Score:     user.Balance,
		Position:  position,
		AvatarURL: user.AvatarURL,
		IsPremium: user.PremiumExpiresAt != nil && user.PremiumExpiresAt.After(time.Now()),
	}
}

func NewUserScoreFromDailyMining(score entity.MiningPerDay, position int64) UserScore {
	return UserScore{
		ID:        score.User.ID,
		Identity:  score.User.Identity(),
		Score:     score.Mined,
		Position:  position,
		AvatarURL: score.User.AvatarURL,
		IsPremium: score.User.IsPremium(),
	}
}

func NewUserScoreFromMonthlyMining(score entity.MiningPerMonth, position int64) UserScore {
	return UserScore{
		ID:        score.User.ID,
		Identity:  score.User.Identity(),
		Score:     score.Mined,
		Position:  position,
		AvatarURL: score.User.AvatarURL,
		IsPremium: score.User.IsPremium(),
	}
}
