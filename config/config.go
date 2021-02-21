package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Init parses the environment variables
func Init() {
	godotenv.Load()
}

// Required checks if all required environment variables are set
func Required() {
	if _, ok := os.LookupEnv("CLIENT_ID"); !ok {
		log.Fatal("CLIENT_ID is missing.")
	}

	if _, ok := os.LookupEnv("CLIENT_SECRET"); !ok {
		log.Fatal("CLIENT_SECRET is missing.")
	}
}

// Get returns the value of an environment variable and a fallback if it's not set
func Get(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	switch key {
	case "PORT":
		return "8080"
	case "AUTHORIZE_PATH":
		return "/authorize"
	case "CALLBACK_PATH":
		return "/callback"
	case "REDIRECT_URI":
		return "http://localhost:8080/callback"
	case "GITHUB_AUTH_ENDPOINT":
		return "https://github.com/login/oauth/authorize"
	case "GITHUB_TOKEN_ENDPOINT":
		return "https://github.com/login/oauth/token"
	case "SCOPE":
		return ""
	default:
		fmt.Printf("%s is an unknown config key", key)
		return ""
	}
}
