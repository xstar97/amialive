package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/xstar97/amialive/internal/config"
)

type JokeResponse struct {
	Error    bool   `json:"error"`
	Type     string `json:"type"`
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
	Joke     string `json:"joke"`
}

// rate limiter
type jokeRateLimiter struct {
	sync.Mutex
	tokens     int
	maxTokens  int
	refillTime time.Duration
}

func newJokeRateLimiter(max int) *jokeRateLimiter {
	rl := &jokeRateLimiter{
		tokens:     max,
		maxTokens:  max,
		refillTime: time.Minute,
	}
	go rl.refill()
	return rl
}

func (rl *jokeRateLimiter) refill() {
	ticker := time.NewTicker(rl.refillTime)
	for range ticker.C {
		rl.Lock()
		rl.tokens = rl.maxTokens
		rl.Unlock()
	}
}

func (rl *jokeRateLimiter) Allow() bool {
	rl.Lock()
	defer rl.Unlock()
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// singleton limiter, initialized when first used
var limiter *jokeRateLimiter
var limiterOnce sync.Once

// GetJoke fetches a random joke from JokeAPI with rate limiting
func GetJoke(cfg *config.Config) string {
	// initialize limiter once
	limiterOnce.Do(func() {
		limiter = newJokeRateLimiter(cfg.JokeRateLimit)
	})

	if !limiter.Allow() {
		return "pong! (JokeAPI rate limit reached, try later 😅)"
	}

	categories := strings.ReplaceAll(cfg.JokeCategory, " ", "")
	url := fmt.Sprintf("https://v2.jokeapi.dev/joke/%s", categories)
	if cfg.JokeSafeMode {
		url += "?safe-mode"
	}

	resp, err := http.Get(url)
	if err != nil {
		return "pong! (couldn't fetch a joke 😢)"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "pong! (error reading joke 😢)"
	}

	var joke JokeResponse
	if err := json.Unmarshal(body, &joke); err != nil {
		return "pong! (bad joke format 🤔)"
	}

	if joke.Type == "single" {
		return fmt.Sprintf("😂 %s", joke.Joke)
	}
	if joke.Type == "twopart" {
		return fmt.Sprintf("😏 %s\n🤣 %s", joke.Setup, joke.Delivery)
	}

	return "pong! (no joke found 😅)"
}
