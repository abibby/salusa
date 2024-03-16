package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/abibby/salusa/env"
	"github.com/abibby/salusa/kernel"
	"github.com/joho/godotenv"
)

var Port int
var DBPath string
var BaseURL string

func Load(ctx context.Context, k *kernel.Kernel) error {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
	} else if err != nil {
		return err
	}

	Port = env.Int("PORT", 2303)
	DBPath = env.String("DATABASE_PATH", "./db.sqlite")
	BaseURL = env.String("BASE_URL", fmt.Sprintf("http://localhost:%d", Port))

	return nil
}
func Kernel() *kernel.KernelConfig {
	return &kernel.KernelConfig{
		Port: Port,
	}
}
