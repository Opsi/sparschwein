package util

import (
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
		return fallback
	}
	return asInt
}
