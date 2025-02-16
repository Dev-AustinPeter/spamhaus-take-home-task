package middleware

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/utils"
)

type RateLimiter struct {
	visitors map[string]time.Time
	mutex    sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]time.Time),
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	log.Println("[INFO] Rate limiter initialized")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mutex.Lock()
		if lastVisit, found := rl.visitors[r.RemoteAddr]; found {
			if time.Since(lastVisit) < 1*time.Second {
				utils.WriteError(w, http.StatusTooManyRequests, fmt.Errorf("%s", "Too many requests"))
				rl.mutex.Unlock()
				return
			}
		}
		rl.visitors[r.RemoteAddr] = time.Now()
		rl.mutex.Unlock()
		next.ServeHTTP(w, r)
	})
}
