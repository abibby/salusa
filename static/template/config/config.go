package config

import (
	"log"

	"github.com/abibby/salusa/env"
	"github.com/joho/godotenv"
)

var Port int
var DBPath string

func Load() error {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Print(err)
	}

	Port = env.Int("PORT", 6900)
	DBPath = env.String("DATABASE_PATH", "./db.sqlite")

	return nil
}
