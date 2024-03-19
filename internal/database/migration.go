package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	observerpkg "github.com/kbgod/coinbot/internal/observer"
	"gorm.io/gorm"
)

type Migrator struct {
	db       *gorm.DB
	observer *observerpkg.Observer
}

func NewMigrator(db *gorm.DB, observer *observerpkg.Observer) *Migrator {
	return &Migrator{
		db:       db,
		observer: observer,
	}
}

func (m *Migrator) RunCommand(cmd string) {
	migrations := m.defineMigrations()
	switch cmd {
	case "migrate":
		if err := migrations.Migrate(); err != nil {
			m.observer.Logger.Fatal().Err(err).Msg("failed to migrate")
		}
	case "rollback":
		if err := migrations.RollbackLast(); err != nil {
			m.observer.Logger.Fatal().Err(err).Msg("failed to rollback")
		}
	case "fresh":
		tables, err := m.db.Migrator().GetTables()
		if err != nil {
			m.observer.Logger.Fatal().Err(err).Msg("get tables")
		}
		for _, t := range tables {
			if err := m.db.Migrator().DropTable(t); err != nil {
				m.observer.Logger.Fatal().Err(err).Str("name", t).Msg("drop table")
			}
		}
		m.observer.Logger.Info().Err(err).Int("count", len(tables)).Msg("dropped all tables")
		if err := migrations.Migrate(); err != nil {
			m.observer.Logger.Fatal().Err(err).Msg("failed to fresh")
		}
	default:
		m.observer.Logger.Fatal().Str("cmd", cmd).Msg("undefined migrator command")
	}
}

func (m *Migrator) defineMigrations() *gormigrate.Gormigrate {
	return gormigrate.New(m.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		initMigration(),
		dailyBoostersMigration(),
		languagesMigration(),
		premiumMigration(),
		leaguesMigration(),
	})
}
