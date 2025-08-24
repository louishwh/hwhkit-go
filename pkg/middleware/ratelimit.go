package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	Rate     int           // 每秒允许的请求数
	Burst    int           // 突发请求数
	Duration time.Duration // 限流窗口时间
	KeyFunc  func(*gin.Context) string // 获取限流键的函数
	ErrorHandler func(*gin.Context) // 限流错误处理函数
}

// DefaultRateLimiterConfig 默认限流配置
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Rate:     100,
		Burst:    10,
		Duration: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded",
			})
			c.Abort()
		},
	}
}

// tokenBucket 令牌桶结构
type tokenBucket struct {
	tokens    int           // 当前令牌数
	capacity  int           // 桶容量
	rate      int           // 令牌生成速率（每秒）
	lastRefill time.Time    // 上次填充时间
	mutex     sync.Mutex   // 互斥锁
}

// newTokenBucket 创建新的令牌桶
func newTokenBucket(capacity, rate int) *tokenBucket {
	return &tokenBucket{
		tokens:    capacity,
		capacity:  capacity,
		rate:      rate,
		lastRefill: time.Now(),
	}
}

// consume 消费令牌
func (tb *tokenBucket) consume() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	
	// 计算应该添加的令牌数
	tokensToAdd := int(elapsed.Seconds()) * tb.rate
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
	
	// 尝试消费一个令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	
	return false
}

// slidingWindow 滑动窗口结构
type slidingWindow struct {
	requests  []time.Time   // 请求时间列表
	limit     int           // 限制数量
	window    time.Duration // 窗口时间
	mutex     sync.Mutex    // 互斥锁
}

// newSlidingWindow 创建新的滑动窗口
func newSlidingWindow(limit int, window time.Duration) *slidingWindow {
	return &slidingWindow{
		requests: make([]time.Time, 0),
		limit:    limit,
		window:   window,
	}
}

// allow 检查是否允许请求
func (sw *slidingWindow) allow() bool {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-sw.window)
	
	// 移除过期的请求
	validRequests := make([]time.Time, 0, len(sw.requests))
	for _, reqTime := range sw.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	sw.requests = validRequests
	
	// 检查是否超过限制
	if len(sw.requests) >= sw.limit {
		return false
	}
	
	// 添加当前请求
	sw.requests = append(sw.requests, now)
	return true
}

// rateLimiter 限流器结构
type rateLimiter struct {
	buckets   map[string]*tokenBucket // 令牌桶映射
	windows   map[string]*slidingWindow // 滑动窗口映射
	config    *RateLimiterConfig
	mutex     sync.RWMutex
}

// newRateLimiter 创建新的限流器
func newRateLimiter(config *RateLimiterConfig) *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string]*tokenBucket),
		windows: make(map[string]*slidingWindow),
		config:  config,
	}
}

// allow 检查是否允许请求
func (rl *rateLimiter) allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	// 使用令牌桶算法
	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = newTokenBucket(rl.config.Burst, rl.config.Rate)
		rl.buckets[key] = bucket
	}
	
	return bucket.consume()
}

// allowSliding 使用滑动窗口检查是否允许请求
func (rl *rateLimiter) allowSliding(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	window, exists := rl.windows[key]
	if !exists {
		window = newSlidingWindow(rl.config.Rate, rl.config.Duration)
		rl.windows[key] = window
	}
	
	return window.allow()
}

// cleanup 清理过期的限流器数据
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			rl.mutex.Lock()
			
			now := time.Now()
			
			// 清理令牌桶（删除超过10分钟未使用的）
			for key, bucket := range rl.buckets {
				if now.Sub(bucket.lastRefill) > 10*time.Minute {
					delete(rl.buckets, key)
				}
			}
			
			// 清理滑动窗口（删除空的窗口）
			for key, window := range rl.windows {
				window.mutex.Lock()
				if len(window.requests) == 0 {
					delete(rl.windows, key)
				}
				window.mutex.Unlock()
			}
			
			rl.mutex.Unlock()
		}
	}()
}

// 全局限流器实例
var globalRateLimiter *rateLimiter

// RateLimit 创建限流中间件
func RateLimit(config ...*RateLimiterConfig) gin.HandlerFunc {
	var cfg *RateLimiterConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	} else {
		cfg = DefaultRateLimiterConfig()
	}
	
	// 创建限流器
	limiter := newRateLimiter(cfg)
	limiter.cleanup() // 启动清理协程
	
	return func(c *gin.Context) {
		key := cfg.KeyFunc(c)
		
		if !limiter.allow(key) {
			cfg.ErrorHandler(c)
			return
		}
		
		c.Next()
	}
}

// RateLimitWithSliding 创建滑动窗口限流中间件
func RateLimitWithSliding(config ...*RateLimiterConfig) gin.HandlerFunc {
	var cfg *RateLimiterConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	} else {
		cfg = DefaultRateLimiterConfig()
	}
	
	limiter := newRateLimiter(cfg)
	limiter.cleanup()
	
	return func(c *gin.Context) {
		key := cfg.KeyFunc(c)
		
		if !limiter.allowSliding(key) {
			cfg.ErrorHandler(c)
			return
		}
		
		c.Next()
	}
}

// RateLimitByIP 基于IP的限流中间件
func RateLimitByIP(rate, burst int) gin.HandlerFunc {
	config := &RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": fmt.Sprintf("Rate limit exceeded: %d requests per second", rate),
			})
			c.Abort()
		},
	}
	
	return RateLimit(config)
}

// RateLimitByUser 基于用户的限流中间件
func RateLimitByUser(rate, burst int) gin.HandlerFunc {
	config := &RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyFunc: func(c *gin.Context) string {
			if userID, exists := GetUserID(c); exists {
				return fmt.Sprintf("user:%d", userID)
			}
			return c.ClientIP() // 回退到IP
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": fmt.Sprintf("Rate limit exceeded: %d requests per second per user", rate),
			})
			c.Abort()
		},
	}
	
	return RateLimit(config)
}

// RateLimitByPath 基于路径的限流中间件
func RateLimitByPath(rate, burst int) gin.HandlerFunc {
	config := &RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyFunc: func(c *gin.Context) string {
			return fmt.Sprintf("%s:%s", c.ClientIP(), c.Request.URL.Path)
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": fmt.Sprintf("Rate limit exceeded for path: %d requests per second", rate),
			})
			c.Abort()
		},
	}
	
	return RateLimit(config)
}

// RateLimitGlobal 全局限流中间件
func RateLimitGlobal(rate, burst int) gin.HandlerFunc {
	config := &RateLimiterConfig{
		Rate:  rate,
		Burst: burst,
		KeyFunc: func(c *gin.Context) string {
			return "global"
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Global rate limit exceeded",
			})
			c.Abort()
		},
	}
	
	return RateLimit(config)
}