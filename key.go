package main

import (
	"log"
	"os"
)

func getEnv(key string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		log.Panicf("%s environment variable is not set!\n", key)
	}
	return env
}
