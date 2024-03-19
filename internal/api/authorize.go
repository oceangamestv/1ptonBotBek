package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kbgod/coinbot/internal/entity"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type userResponse struct {
	ID                      int64      `json:"id"`
	Balance                 int64      `json:"balance"`
	FirstName               string     `json:"first_name"`
	Energy                  int32      `json:"energy"`
	EnergyLevel             int32      `json:"energy_level"`
	MaxEnergyLevel          int32      `json:"max_energy_level"`
	MineLevel               int32      `json:"mine_level"`
	AutoFarmer              bool       `json:"auto_farmer"`
	AutoFarmerProfit        int64      `json:"auto_farmer_profit"`
	AccessToken             *string    `json:"access_token"`
	AccessTokenExpiresAt    *time.Time `json:"access_token_expires_at"`
	PremiumExpiresAt        *time.Time `json:"premium_expires_at"`
	IsPremium               bool       `json:"is_premium"`
	DailyBoosterAvailableAt time.Time  `json:"daily_booster_available_at"`
}

func newUserResponse(u *entity.User) *userResponse {
	return &userResponse{
		ID:                      u.ID,
		Balance:                 u.Balance,
		FirstName:               u.FirstName,
		Energy:                  u.CurrentEnergy(),
		EnergyLevel:             u.EnergyLevel,
		MaxEnergyLevel:          u.MaxEnergyLevel,
		MineLevel:               u.MineLevel,
		AutoFarmer:              u.AutoFarmer,
		AutoFarmerProfit:        u.AutoFarmerProfit,
		AccessToken:             u.WebAppAccessToken,
		AccessTokenExpiresAt:    u.WebAppAccessTokenExpiresAt,
		PremiumExpiresAt:        u.PremiumExpiresAt,
		IsPremium:               u.PremiumExpiresAt != nil && u.PremiumExpiresAt.After(time.Now()),
		DailyBoosterAvailableAt: u.DailyBoosterAvailableAt,
	}
}

func (h *handler) getMe(ctx *fiber.Ctx) error {
	user := getUserFromCtx(ctx)
	return ctx.JSON(newUserResponse(user))
}

func (h *handler) authorize(ctx *fiber.Ctx) error {
	//req := map[string]string{}
	req, err := url.ParseQuery(string(ctx.Body()))
	if err != nil {
		return newMessageResponse(ctx, http.StatusBadRequest, err.Error())
	}

	hash := req.Get("hash")
	if hash == "" {
		return newMessageResponse(ctx, http.StatusBadRequest, "hash is required")
	}
	delete(req, "hash")
	fmt.Println(hash)

	data := make(map[string]string, len(req))
	for k := range req {
		data[k] = fmt.Sprint(req.Get(k))
	}

	fmt.Println(data)

	authUser, err := h.validateAndExtractTelegramUserData(data, h.svc.CFG.BotToken, hash)
	if err != nil {
		return newMessageResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if err := h.svc.AuthorizeByWebApp(authUser); err != nil {
		return err
	}

	return ctx.JSON(newUserResponse(authUser))
}

func (h *handler) validateAndExtractTelegramUserData(
	authData map[string]string, secret, signature string,
) (*entity.User, error) {
	dataCheckArr := make([]string, 0, len(authData))
	for key, value := range authData {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, value))
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	fmt.Println(dataCheckString)

	mac := hmac.New(sha256.New, []byte("WebAppData"))
	mac.Write([]byte(secret))

	mac = hmac.New(sha256.New, mac.Sum(nil))
	mac.Write([]byte(dataCheckString))
	calculatedSignature := hex.EncodeToString(mac.Sum(nil))
	if calculatedSignature != signature {
		return nil, errors.New("invalid auth signature, given: " + signature)
	}
	if userData, ok := authData["user"]; ok {
		userFields := map[string]any{}
		if err := json.Unmarshal([]byte(userData), &userFields); err != nil {
			return nil, errors.New("invalid auth data")
		}
		for k, v := range userFields {
			if floatval, ok := v.(float64); ok {
				authData[k] = strconv.FormatFloat(floatval, 'f', -1, 64)
				continue
			}
			authData[k] = fmt.Sprint(v)
		}
	}
	telegramID, err := strconv.ParseInt(authData["id"], 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth data")
	}

	authDate, err := strconv.ParseInt(authData["auth_date"], 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth data")
	}
	if (time.Now().Unix() - authDate) > 86400 {
		return nil, errors.New("auth token expired")
	}
	u := &entity.User{
		ID: telegramID,
	}
	if authData["username"] != "" {
		username := authData["username"]
		u.Username = username
	}
	if authData["first_name"] != "" {
		firstName := authData["first_name"]
		u.FirstName = firstName
	}
	if authData["avatar"] != "" {
		avatar := authData["avatar"]
		u.AvatarURL = &avatar
	}
	fmt.Println(authData)

	return u, nil
}
