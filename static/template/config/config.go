package config

import (
	"errors"
	"os"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	DBPath   string
	Database dialects.Config
	Mail     email.Config
}

func Load() *Config {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
		// fall through
	} else if err != nil {
		panic(err)
	}

	return &Config{
		Port:     env.Int("PORT", 2303),
		Database: sqlite.NewConfig(env.String("DATABASE_PATH", "./db.sqlite")),
		Mail: &email.SMTPConfig{
			From:     env.String("MAIL_FROM", "salusa@example.com"),
			Host:     env.String("MAIL_HOST", "sandbox.smtp.mailtrap.io"),
			Port:     env.Int("MAIL_PORT", 2525),
			Username: env.String("MAIL_USERNAME", "user"),
			Password: env.String("MAIL_PASSWORD", "pass"),
		},
	}
}

func (c *Config) GetHTTPPort() int {
	return c.Port
}
