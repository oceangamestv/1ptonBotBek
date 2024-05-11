package main

import (
	"context"
	"fmt"
	"github.com/kbgod/coinbot/config"
	"github.com/kbgod/coinbot/internal/api"
	"github.com/kbgod/coinbot/internal/database"
	"github.com/kbgod/coinbot/internal/handler"
	observerpkg "github.com/kbgod/coinbot/internal/observer"
	"github.com/kbgod/coinbot/internal/service"
	"github.com/kbgod/illuminate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const svcTemplate = `
[Unit]
Description=Coinbot Telegram Bot
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory={path}
EnvironmentFile={path}/.env
ExecStart={path}/bot migrate
Restart=always

[Install]
WantedBy=multi-user.target`

func setupSystemCtlService() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	tmpl := strings.ReplaceAll(svcTemplate, "{path}", exPath)
	fmt.Println(tmpl)
	f, err := os.Create("/etc/systemd/system/coinbot.service")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(tmpl)
	if err != nil {
		panic(err)
	}
}
func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	observer := observerpkg.New(cfg.LogLevel, cfg.Debug)

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
	if err != nil {
		observer.Logger.Fatal().Err(err).Msg("connect to database")
	}
	if cfg.DBDebug {
		db = db.Debug()
	}

	migrator := database.NewMigrator(db, observer)

	if len(os.Args) > 1 {
		if os.Args[1] == "setup" {
			setupSystemCtlService()
			return
		}
		if os.Args[1] == "fresh" && !cfg.FreshAllowed {
			observer.Logger.Fatal().Msg("fresh command not allowed")
		}
		migrator.RunCommand(os.Args[1])
	}

	svc := service.New(cfg, db, observer)

	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-exit
		cancel()
	}()

	go api.Run(svc)

	botClient, err := illuminate.NewBot(cfg.BotToken, nil)
	if err != nil {
		observer.Logger.Fatal().Err(err).Msg("create bot client")
	}
	observer.Logger.Info().Str("username", botClient.Username).Msg("bot authorized")

	h := handler.New(svc, botClient, botClient.User)

	if err := h.Run(ctx); err != nil {
		observer.Logger.Fatal().Err(err).Msg("run handler")
	}
}
