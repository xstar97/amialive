package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xstar97/amialive/internal/api"
)

// Config holds runtime configuration
type Config struct {
	Port           string
	JokeChance     int    // percentage chance (0–100)
	JokeCategory   string // Valid JokeAPI category string
	JokeSafeMode   bool   // true if safe-mode is enabled
	JokeRateLimit  int    // max requests per minute (1–120)
}

// LoadConfig loads environment variables with defaults
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	jokeChance := getEnvInt("JOKE_CHANCE", 30)

	rawCategory := getEnv("JOKE_CATEGORY", "Programming")
	jokeCategory := validateJokeCategory(rawCategory)

	jokeSafeMode := getEnvBool("JOKE_SAFEMODE", true)
	jokeRateLimit := getEnvInt("JOKES_REQUESTS", 60)

	// Clamp to valid range 1–120
	if jokeRateLimit < 1 {
		jokeRateLimit = 1
	}
	if jokeRateLimit > 120 {
		jokeRateLimit = 120
	}

	return &Config{
		Port:          port,
		JokeChance:    jokeChance,
		JokeCategory:  jokeCategory,
		JokeSafeMode:  jokeSafeMode,
		JokeRateLimit: jokeRateLimit,
	}
}

// validateJokeCategory filters out invalid categories using api.ValidJokeCategories
func validateJokeCategory(raw string) string {
	parts := strings.Split(raw, ",")
	validParts := []string{}

	for _, p := range parts {
		p = strings.TrimSpace(p)
		for _, valid := range api.ValidJokeCategories {
			if strings.EqualFold(p, valid) {
				validParts = append(validParts, valid)
				break
			}
		}
	}

	if len(validParts) == 0 {
		log.Printf("No valid JOKE_CATEGORY found, defaulting to Any")
		return "Any"
	}

	return strings.Join(validParts, ",")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Invalid value for %s, using default %d\n", key, fallback)
		return fallback
	}
	return num
}

func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	switch strings.ToLower(val) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		log.Printf("Invalid value for %s, using default %v\n", key, fallback)
		return fallback
	}
}
