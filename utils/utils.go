package utils

import (
	"os"
)

// GetEnvString returns a string from the provided environment variable
func GetEnvString(envVar string, defaults string) string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaults
	}
	return value
}
