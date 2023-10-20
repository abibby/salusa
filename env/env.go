package env

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func String(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}

func Bool(key string, defaultValue bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	v = strings.ToLower(v)
	return v == "true" || v == "1"
}

func Int(key string, defaultValue int) int {
	vStr, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	v, err := strconv.Atoi(vStr)
	if err != nil {
		slog.Warn("failed to parse int", slog.String("value", vStr), slog.Any("error", err))
		return defaultValue
	}
	return v
}

func Float64(key string, defaultValue float64) float64 {
	vStr, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	v, err := strconv.ParseFloat(vStr, 64)
	if err != nil {
		slog.Warn("failed to parse float", slog.String("value", vStr), slog.Any("error", err))
		return defaultValue
	}
	return v
}

func Float32(key string, defaultValue float32) float32 {
	return float32(Float64(key, float64(defaultValue)))
}
