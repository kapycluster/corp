package util

import (
	"log"
	"os"
)

func MustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env var not set: %s", key)
	}

	return v
}
