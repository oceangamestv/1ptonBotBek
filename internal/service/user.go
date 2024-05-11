package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/kbgod/coinbot/internal/entity"
	"github.com/kbgod/coinbot/internal/valueobject"
	"github.com/kbgod/coinbot/pkg/randomizer"
	"github.com/kbgod/illuminate"
	"gorm.io/gorm"
)

type GetUserOptions struct {
	TgUser    *illuminate.User
	IsPrivate bool
	Promo     *string
	AvatarURL *string
}

func (s *Service) GetUser(opts *GetUserOptions) (*entity.User, error) {
	var user entity.User
	if err := s.db.Take(&user, opts.TgUser.Id).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get user: %w", err)
	} else if err == nil {
		mustUpdate := make(map[string]any)
		if opts.IsPrivate && user.StoppedAt != nil {
			mustUpdate["stopped_at"] = nil
		}
		if opts.TgUser.FirstName != user.FirstName {
			mustUpdate["first_name"] = opts.TgUser.FirstName
		}
		if opts.TgUser.Username != user.Username {
			mustUpdate["username"] = opts.TgUser.Username
		}
		if opts.TgUser.LanguageCode != "" && opts.TgUser.LanguageCode != user.LanguageCode {
			mustUpdate["language_code"] = opts.TgUser.LanguageCode
		}

		if len(mustUpdate) > 0 {
			if err := s.db.Model(&entity.User{}).Where("id", user.ID).Updates(mustUpdate).Error; err != nil {
				return nil, fmt.Errorf("update user: %w", err)
			}
		}

		return &user, nil
	}

	user.ID = opts.TgUser.Id
	user.FirstName = opts.TgUser.FirstName
	user.Username = opts.TgUser.Username
	user.LanguageCode = opts.TgUser.LanguageCode

	if opts.Promo != nil {
		if strings.HasPrefix(*opts.Promo, "r_") {
			refererID, err := strconv.ParseInt(strings.TrimPrefix(*opts.Promo, "r_"), 10, 64)
			if err == nil {
				user.RefererID = &refererID
			}
			if err := s.db.Model(&entity.User{}).Where("id = ?", refererID).Updates(map[string]any{
				"balance_usd":         gorm.Expr("balance_usd + ?", ReferralRegisterReward),
				"referral_profit_usd": gorm.Expr("referral_profit_usd + ?", ReferralRegisterReward),
			}).Error; err != nil {
				return nil, fmt.Errorf("update referer balance: %w", err)
			}
		} else {
			p, err := entity.GetPromoByName(s.db, *opts.Promo)
			if err != nil {
				return nil, fmt.Errorf("get promo by name: %w", err)
			}
			user.PromoID = &p.ID
			user.Balance = p.Charge
		}
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (s *Service) SetUserBotState(user *entity.User, state string, stateContext ...string) error {
	mustUpdate := make(map[string]any, 2)
	if state == "" {
		mustUpdate["bot_state"] = nil
		mustUpdate["bot_state_context"] = nil
	} else {
		mustUpdate["bot_state"] = state
		if len(stateContext) > 0 {
			mustUpdate["bot_state_context"] = stateContext[0]
		}
	}
	return s.db.Model(&entity.User{}).Where("id", user.ID).Updates(mustUpdate).Error
}

func (s *Service) RemoveUserBotState(user *entity.User) error {
	return s.db.Model(&entity.User{}).Where("id", user.ID).Updates(map[string]any{
		"bot_state":         nil,
		"bot_state_context": nil,
	}).Error
}
func (s *Service) ProcessAutoFarmer(u *entity.User) error {
	if u.AutoFarmer && u.LastMineAt.Add(1*time.Hour).Before(time.Now()) {
		secondsFromLastMining := time.Now().Sub(u.LastMineAt).Seconds()
		// if bigger than 12 hours, reduce to 12 hours
		if secondsFromLastMining > 12*60*60 {
			secondsFromLastMining = 12 * 60 * 60
		}
		mined := int64(int32(secondsFromLastMining) * u.EnergyLevel)
		if mined > 0 {
			if err := s.db.Model(&entity.User{}).Where("id = ?", u.ID).
				Updates(map[string]any{
					"balance":      gorm.Expr("balance + ?", mined),
					"energy":       u.MaxEnergy(),
					"last_mine_at": time.Now(),
				}).Error; err != nil {
				return fmt.Errorf("update user balance: %w", err)
			}
		}
		u.LastMineAt = time.Now()
		u.Balance += mined
		u.AutoFarmerProfit = mined
		u.Energy = u.MaxEnergy()
	}

	return nil
}
func (s *Service) AuthorizeByWebApp(user *entity.User) error {
	var avatarURL *string
	if user.AvatarURL != nil {
		avatarURL = new(string)
		*avatarURL = *user.AvatarURL
	}
	//if err := s.db.Take(user, user.ID).Error; err != nil {
	//	return fmt.Errorf("get user: %w", err)
	//}

	u, err := s.GetUser(&GetUserOptions{
		TgUser: &illuminate.User{
			Id:        user.ID,
			FirstName: user.FirstName,
			Username:  user.Username,
		},
		IsPrivate: true,
	})
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if err := s.ProcessAutoFarmer(u); err != nil {
		return fmt.Errorf("process auto farmer: %w", err)
	}
	*user = *u

	accessToken, err := randomizer.GenerateRandomHex(64)
	if err != nil {
		return fmt.Errorf("generate random hex: %w", err)
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	mustUpdate := map[string]any{
		"web_app_access_token":            accessToken,
		"web_app_access_token_expires_at": expiresAt,
		"avatar_url":                      avatarURL,
	}
	return s.db.Model(user).Updates(mustUpdate).Error
}

func (s *Service) GetUserByAccessToken(accessToken string) (*entity.User, error) {
	var user entity.User
	if err := s.db.Where("web_app_access_token = ?", accessToken).Take(&user).Error; err != nil {
		return nil, fmt.Errorf("get user by access token: %w", err)
	}
	if err := s.ProcessAutoFarmer(&user); err != nil {
		return nil, fmt.Errorf("process auto farmer: %w", err)
	}
	return &user, nil
}

type Leaderboard struct {
	Players []valueobject.UserScore `json:"players"`
	Me      *valueobject.UserScore  `json:"me"`
}

func (s *Service) Leaderboard(u *entity.User) (*Leaderboard, error) {
	var users []entity.User
	if err := s.db.
		Order("balance desc").
		Where("balance > 0").
		Limit(100).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("get leaderboard: %w", err)
	}
	var publicUsers []valueobject.UserScore
	var me *valueobject.UserScore
	for idx, user := range users {
		score := valueobject.NewUserScoreFromEntity(user, int64(idx+1))
		publicUsers = append(publicUsers, score)

		if user.ID == u.ID {
			me = &score
		}
	}
	if me == nil {
		var position int64
		if err := s.db.
			Model(&entity.User{}).
			Where("balance > 0 and balance > ?", u.Balance).
			Limit(100).
			Count(&position).Error; err != nil {
			return nil, fmt.Errorf("get user position: %w", err)
		}

		score := valueobject.NewUserScoreFromEntity(*u, position+1)
		me = &score
	}

	return &Leaderboard{
		Players: publicUsers,
		Me:      me,
	}, nil
}

func (s *Service) DailyLeaderboard(u *entity.User) (*Leaderboard, error) {
	var scores []entity.MiningPerDay
	if err := s.db.
		Order("mined desc").
		Joins("User").
		Where("date = ?", now.With(time.Now().UTC()).BeginningOfDay()).Limit(100).
		Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("get daily leaderboard: %w", err)
	}
	var publicScores []valueobject.UserScore
	var me *valueobject.UserScore
	for idx, score := range scores {
		publicScore := valueobject.NewUserScoreFromDailyMining(score, int64(idx+1))
		publicScores = append(publicScores, publicScore)

		if score.UserID == u.ID {
			me = &publicScore
		}
	}
	if me == nil {
		var meScore entity.MiningPerDay
		if err := s.db.
			Where("date = ? and user_id = ?", now.With(time.Now().UTC()).BeginningOfDay(), u.ID).
			Take(&meScore).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user score: %w", err)
		} else if err == nil {
			var position int64
			if err := s.db.
				Model(&entity.MiningPerDay{}).
				Where("date = ? and mined > ?", now.With(time.Now().UTC()).BeginningOfDay(), meScore.Mined).
				Count(&position).Error; err != nil {
				return nil, fmt.Errorf("get user position: %w", err)
			}
			meScore.User = u
			score := valueobject.NewUserScoreFromDailyMining(meScore, position+1)
			me = &score
		}
	}

	return &Leaderboard{
		Players: publicScores,
		Me:      me,
	}, nil
}

func (s *Service) MonthlyLeaderboard(u *entity.User) (*Leaderboard, error) {
	var scores []entity.MiningPerMonth
	if err := s.db.
		Order("mined desc").
		Joins("User").
		Where("date = ?", now.With(time.Now().UTC()).BeginningOfMonth()).Limit(100).
		Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("get monthly leaderboard: %w", err)
	}
	var publicScores []valueobject.UserScore
	var me *valueobject.UserScore
	for idx, score := range scores {
		publicScore := valueobject.NewUserScoreFromMonthlyMining(score, int64(idx+1))
		publicScores = append(publicScores, publicScore)

		if score.UserID == u.ID {
			me = &publicScore
		}
	}
	if me == nil {
		var meScore entity.MiningPerMonth
		if err := s.db.
			Where("date = ? and user_id = ?", now.With(time.Now().UTC()).BeginningOfMonth(), u.ID).
			Take(&meScore).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user score: %w", err)
		} else if err == nil {
			var position int64
			if err := s.db.
				Model(&entity.MiningPerMonth{}).
				Where("date = ? and mined > ?", now.With(time.Now().UTC()).BeginningOfMonth(), meScore.Mined).
				Count(&position).Error; err != nil {
				return nil, fmt.Errorf("get user position: %w", err)
			}
			meScore.User = u
			score := valueobject.NewUserScoreFromMonthlyMining(meScore, position+1)
			me = &score
		}
	}

	return &Leaderboard{
		Players: publicScores,
		Me:      me,
	}, nil
}

func (s *Service) GetReferralsCount(user *entity.User) (int64, error) {
	var count int64
	if err := s.db.Model(&entity.User{}).Where("referer_id = ?", user.ID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("get referrals count: %w", err)
	}
	return count, nil
}

func (s *Service) UpdateUserStoppedStatus(user *entity.User, stopped bool) error {
	var stoppedAt *time.Time
	if stopped {
		stoppedAt = new(time.Time)
		*stoppedAt = time.Now()
	}
	return s.db.Model(&entity.User{}).Where("id", user.ID).Update("stopped_at", stoppedAt).Error
}

func (s *Service) BatchActiveUsers(size int, callback func(entity.User)) error {
	var users []entity.User
	return s.db.Model(&entity.User{}).
		Where("stopped_at is null").
		Where("language_code != 'uk'").
		Order("id").
		Limit(size).
		FindInBatches(&users, size, func(tx *gorm.DB, batch int) error {
			for _, user := range users {
				callback(user)
			}
			return nil
		}).Error

}
