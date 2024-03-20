package config

import (
	"errors"
	"os"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	DBPath   string
	Database dialects.Config
}

func Load() *Config {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
		// fall through
	} else if err != nil {
		panic(err)
	}

	Port := env.Int("PORT", 2303)

	return &Config{
		Port:     Port,
		Database: sqlite.NewConfig(env.String("DATABASE_PATH", "./db.sqlite")),
		// Database: &databasedi.SimpleMySQLConfig{
		// 	Username: "",
		// 	Password: "",
		// 	Address:  "",
		// 	Database: "",
		// },
	}
}

func (c *Config) GetHTTPPort() int {
	return c.Port
}
