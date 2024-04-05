package config

import (
	"errors"
	"os"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/env"
	"github.com/abibby/salusa/event"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	BasePath string

	Database dialects.Config
	Mail     email.Config
	Queue    event.Config
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
		BasePath: env.String("BASE_PATH", ""),
		Database: sqlite.NewConfig(env.String("DATABASE_PATH", "./db.sqlite")),
		Queue:    event.NewChannelQueueConfig(),
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
func (c *Config) GetBaseURL() string {
	return c.BasePath
}

func (c *Config) DBConfig() dialects.Config {
	return c.Database
}
func (c *Config) MailConfig() email.Config {
	return c.Mail
}
func (c *Config) QueueConfig() event.Config {
	return c.Queue
}
