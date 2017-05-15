package main

import "os"

func getEnv(key string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		panic("ATLAS_USERNAME environment variable is not set!\n")
	}
	return env
}
