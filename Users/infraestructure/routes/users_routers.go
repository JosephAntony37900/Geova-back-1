package routes

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/controllers"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter gestiona los limitadores por IP con expiración automática
type RateLimiter struct {
	ips             map[string]*limiterEntry
	mu              sync.RWMutex
	r               rate.Limit
	b               int
	ttl             time.Duration
	cleanupInterval time.Duration
}

// RateLimiterConfig contiene la configuración del rate limiter
type RateLimiterConfig struct {
	RequestsPerSecond float64
	Burst             int
	TTL               time.Duration
	CleanupInterval   time.Duration
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		ips:             make(map[string]*limiterEntry),
		r:               rate.Limit(config.RequestsPerSecond),
		b:               config.Burst,
		ttl:             config.TTL,
		cleanupInterval: config.CleanupInterval,
	}

	go rl.cleanupExpiredEntries()

	return rl
}

// GetLimiter obtiene o crea un limiter para una IP (optimizado con RLock/Lock)
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	// Intento 1: lectura rápida con RLock (solo lectura, sin modificar lastSeen)
	rl.mu.RLock()
	entry, exists := rl.ips[ip]
	if exists {
		limiter := entry.limiter
		rl.mu.RUnlock()
		return limiter
	}
	rl.mu.RUnlock()

	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists = rl.ips[ip]
	if exists {
		entry.lastSeen = now
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.r, rl.b)
	rl.ips[ip] = &limiterEntry{
		limiter:  limiter,
		lastSeen: now,
	}

	return limiter
}

func (rl *RateLimiter) cleanupExpiredEntries() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, entry := range rl.ips {
			if now.Sub(entry.lastSeen) > rl.ttl {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)
		
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Demasiadas peticiones. Intenta más tarde.",
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// getEnvFloat obtiene un float desde variable de entorno o usa default
func getEnvFloat(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

// getEnvInt obtiene un int desde variable de entorno o usa default
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// getEnvDuration obtiene una duración desde variable de entorno o usa default
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func SetupUserRoutes(r *gin.Engine, 
	createUserController *controllers.CreateUserController,
	getUsersController *controllers.GetAllUsersController,
	getUsersControllerById *controllers.GetUserByIdController,
	updateUserController *controllers.UpdateUserController,
	deleteUserController *controllers.DeleteUserController,
	loginUserController *controllers.LoginUserController,
) {
	
	loginLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("USERS_LOGIN_RATE_LIMIT", 4),
		Burst:             getEnvInt("USERS_LOGIN_BURST_LIMIT", 3),
		TTL:               getEnvDuration("USERS_RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("USERS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})
	
	registerLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("USERS_REGISTER_RATE_LIMIT", 0.5),
		Burst:             getEnvInt("USERS_REGISTER_BURST_LIMIT", 2),
		TTL:               getEnvDuration("USERS_RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("USERS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})
	
	modifyLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("USERS_MODIFY_RATE_LIMIT", 3),
		Burst:             getEnvInt("USERS_MODIFY_BURST_LIMIT", 5),
		TTL:               getEnvDuration("USERS_RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("USERS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})
	
	readLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("USERS_READ_RATE_LIMIT", 20),
		Burst:             getEnvInt("USERS_READ_BURST_LIMIT", 40),
		TTL:               getEnvDuration("USERS_RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("USERS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})

	loginRoutes := r.Group("/users")
	loginRoutes.Use(loginLimiter.RateLimitMiddleware())
	{
		loginRoutes.POST("/login", loginUserController.Execute)
	}

	registerRoutes := r.Group("/users")
	registerRoutes.Use(registerLimiter.RateLimitMiddleware())
	{
		registerRoutes.POST("", createUserController.Execute)
	}

	modifyRoutes := r.Group("/users")
	modifyRoutes.Use(modifyLimiter.RateLimitMiddleware())
	{
		modifyRoutes.PUT("/:id", updateUserController.Execute)
		modifyRoutes.DELETE("/:id", deleteUserController.Execute)
	}

	readRoutes := r.Group("/users")
	readRoutes.Use(readLimiter.RateLimitMiddleware())
	{
		readRoutes.GET("", getUsersController.Execute)
		readRoutes.GET("/:id", getUsersControllerById.Execute)
	}
}