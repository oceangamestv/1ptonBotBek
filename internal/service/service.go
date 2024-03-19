package service

import (
	"github.com/kbgod/coinbot/config"
	observerpkg "github.com/kbgod/coinbot/internal/observer"
	"gorm.io/gorm"
)

type Service struct {
	CFG      *config.Config
	db       *gorm.DB
	Observer *observerpkg.Observer
}

func New(CFG *config.Config, db *gorm.DB, observer *observerpkg.Observer) *Service {
	return &Service{CFG: CFG, db: db, Observer: observer}
}
