package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/you/sharing-vision-backend-v2/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu    sync.Mutex
	store map[string]*rate.Limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{store: make(map[string]*rate.Limiter)}
}

func (rl *RateLimiter) getLimiter(ip string, rps int, burst int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, ok := rl.store[ip]; ok {
		return limiter
	}
	limiter := rate.NewLimiter(rate.Every(time.Duration(1_000_000_000/rps)*time.Nanosecond), burst)
	rl.store[ip] = limiter
	return limiter
}

func RateLimitPublic(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.getLimiter(ip, config.Conf.RateLimit.PublicRPS, config.Conf.RateLimit.PublicBurst).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RateLimitAdmin(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.getLimiter(ip, config.Conf.RateLimit.AdminRPS, config.Conf.RateLimit.AdminBurst).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "admin rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		config.Log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.String("latency", latency.String()),
		)
	}
}
