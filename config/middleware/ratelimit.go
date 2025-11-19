package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// GlobalRateLimiter returns a middleware that limits the total incoming request rate.
// If rps <= 0 the limiter is disabled (no throttling).
func GlobalRateLimiter(rps float64, burst int, skip map[string]struct{}) gin.HandlerFunc {
	if rps <= 0 { // disabled
		return func(c *gin.Context) { c.Next() }
	}
	if burst <= 0 { // ensure sane burst
		burst = int(rps)
		if burst < 1 {
			burst = 1
		}
	}
	lim := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		// Skip OPTIONS preflight quickly
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		// Skip specific full paths if configured
		if skip != nil {
			if _, ok := skip[c.FullPath()]; ok {
				c.Next()
				return
			}
		}
		if !lim.Allow() {
			// Optional small retry-after hint (approximate: next token)
			retryAfter := time.Until(time.Now().Add(lim.Reserve().Delay())) / time.Millisecond
			c.Header("Retry-After", "1") // seconds (coarse)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":          "rate limit exceeded",
				"retry_after_ms": retryAfter,
			})
			return
		}
		c.Next()
	}
}
