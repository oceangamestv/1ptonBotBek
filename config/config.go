package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     string `env:"POSTGRES_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
	FreshAllowed     bool   `env:"FRESH_ALLOWED" envDefault:"false"`

	API struct {
		Host string `env:"API_HOST" envDefault:"localhost"`
		Port string `env:"API_PORT" envDefault:"8080"`
		TLS  bool   `env:"API_TLS" envDefault:"false"`
	}

	BotToken string `env:"BOT_TOKEN,required"`

	LogLevel    zerolog.Level `env:"LOG_LEVEL" envDefault:"debug"`
	Debug       bool          `env:"DEBUG" envDefault:"false"`
	DBDebug     bool          `env:"DB_DEBUG" envDefault:"false"`
	FrontendURL string        `env:"FRONTEND_URL" envDefault:"http://localhost:3000"`

	Name        string `env:"NAME" envDefault:"coinbot"`
	NewsChannel string `env:"NEWS_CHANNEL" envDefault:"coinbot_news"`
}

func New() (*Config, error) {
	_ = godotenv.Load(".env")
	f := &Config{}
	return f, env.Parse(f)
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
	)
}
