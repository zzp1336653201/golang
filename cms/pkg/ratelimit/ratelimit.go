package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"cms/internal/model"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶限流
type TokenBucket struct {
	rate       int           // 每秒补充的令牌数
	bucket     int           // 桶容量
	tokens     int           // 当前令牌数
	lastUpdate time.Time     // 上次更新时间
	mu         sync.Mutex
}

func NewTokenBucket(rate, burst int) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		bucket:     burst,
		tokens:     burst,
		lastUpdate: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	now := time.Now()
	// 补充令牌
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.tokens += int(elapsed * float64(tb.rate))
	if tb.tokens > tb.bucket {
		tb.tokens = tb.bucket
	}
	tb.lastUpdate = now
	
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(rateLimitCfg struct {
	QPS   int
	Burst int
}) gin.HandlerFunc {
	bucket := NewTokenBucket(rateLimitCfg.QPS, rateLimitCfg.Burst)
	
	return func(c *gin.Context) {
		if !bucket.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, model.Error("请求过于频繁，请稍后重试"))
			return
		}
		c.Next()
	}
}
