package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// ValidJokeCategories is the single source of truth
var ValidJokeCategories = []string{
	"Programming", "Miscellaneous", "Dark", "Pun", "Spooky", "Christmas", "Any",
}

// Config holds runtime configuration
type Config struct {
	Port           string
	JokeChance     int    // percentage chance (0–100)
	JokeCategory   string // Valid JokeAPI category string
	JokeSafeMode   bool   // true if safe-mode is enabled
	JokeRateLimit  int    // max requests per minute (1–120)
}

// Load loads configuration from flags, environment variables, and defaults
func Load() *Config {
	cfg := &Config{}

	// Define flags
	flag.StringVar(&cfg.Port, "port", "", "Port to run the server on")
	flag.IntVar(&cfg.JokeChance, "joke-chance", -1, "Percentage chance to show a joke (0-100)")
	flag.StringVar(&cfg.JokeCategory, "joke-category", "", "Comma-separated JokeAPI categories")
	flag.BoolVar(&cfg.JokeSafeMode, "joke-safe-mode", true, "Enable safe-mode for jokes")
	flag.IntVar(&cfg.JokeRateLimit, "jokes-requests", -1, "Max JokeAPI requests per minute (1-120)")

	flag.Parse()

	// Port
	if cfg.Port == "" {
		cfg.Port = getEnv("PORT", "8080")
	}

	// Joke chance
	if cfg.JokeChance == -1 {
		cfg.JokeChance = getEnvInt("JOKE_CHANCE", 30)
	}

	// Joke category
	if cfg.JokeCategory == "" {
		cfg.JokeCategory = getEnv("JOKE_CATEGORY", "Programming")
	}
	cfg.JokeCategory = validateJokeCategory(cfg.JokeCategory)

	// Safe mode
	cfg.JokeSafeMode = getEnvBool("JOKE_SAFEMODE", cfg.JokeSafeMode)

	// Joke rate limit
	if cfg.JokeRateLimit == -1 {
		cfg.JokeRateLimit = getEnvInt("JOKES_REQUESTS", 60)
	}
	if cfg.JokeRateLimit < 1 {
		cfg.JokeRateLimit = 1
	}
	if cfg.JokeRateLimit > 120 {
		cfg.JokeRateLimit = 120
	}

	return cfg
}

// validateJokeCategory filters out invalid categories using ValidJokeCategories
func validateJokeCategory(raw string) string {
	parts := strings.Split(raw, ",")
	validParts := []string{}

	for _, p := range parts {
		p = strings.TrimSpace(p)
		for _, valid := range ValidJokeCategories {
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

// helper functions
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
