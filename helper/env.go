package helper

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Get env variable from file using godotenv library
func GetEnv(key, defaultValue string) string {
	err := godotenv.Load()
	CheckError(err, "Error loading .env file")
	value := os.Getenv(key)
	if len(value) == 0 {
		fmt.Println(defaultValue)
		return defaultValue
	}
	return value
}