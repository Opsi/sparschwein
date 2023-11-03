package util

import (
	"log/slog"
	"os"
	"strconv"
)

func LookupStringEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func LookupIntEnv(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	asInt, err := strconv.Atoi(value)
	if err != nil {
		slog.Warn("Could not parse environment variable %s as int: %s", key, err)
		return fallback
	}
	return asInt
}
