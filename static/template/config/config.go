package config

import (
	"context"
	"errors"
	"os"

	"github.com/abibby/salusa/env"
	"github.com/abibby/salusa/kernel"
	"github.com/joho/godotenv"
)

var Port int
var DBPath string

func Load(ctx context.Context) error {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
	} else if err != nil {
		return err
	}

	Port = env.Int("PORT", 6900)
	DBPath = env.String("DATABASE_PATH", "./db.sqlite")

	return nil
}
func Kernel() *kernel.KernelConfig {
	return &kernel.KernelConfig{
		Port: Port,
	}
}
