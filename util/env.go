package util

import (
	"log"
	"os"
)

func GetOrCrash(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing env variable %s\n", key)
	}
	return value
}

func GetOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Missing env variable %s using default %s\n", key, defaultValue)
		return defaultValue
	}
	log.Printf("Have env variable %s using %s\n", key, defaultValue)
	return value
}
