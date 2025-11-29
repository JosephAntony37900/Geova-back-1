package routes

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
)

// limiterEntry almacena un rate limiter con su timestamp de último uso
type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter gestiona los limitadores por IP con expiración automática
type RateLimiter struct {
	ips       map[string]*limiterEntry
	mu        sync.RWMutex
	r         rate.Limit
	b         int
	ttl       time.Duration
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

	// Iniciar goroutine de limpieza
	go rl.cleanupExpiredEntries()

	return rl
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
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

// cleanupExpiredEntries elimina periódicamente las entradas antiguas
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

func SetUpProjectsRoutes(r *gin.Engine, 
	createProjectController *controllers.CreateProjectController,
	getProjectsController *controllers.GetAllProjectsController,
	getProjectByIdController *controllers.GetProjectByIdController,
	getProjectByNameController *controllers.GetProjectByNameController,
	getProjectByCategoryController *controllers.GetProjectByCategoryController,
	getProjectByDateController *controllers.GetProjectByDateController,
	getProjectsStats *controllers.GetProjectStatsController,
	updateProjectController *controllers.UpdateProjectController,
	deleteProjectController *controllers.DeleteProjectController,
	getProjectByUserId *controllers.GetProjectsByUserIdController,
) {

	writeLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("PROJECTS_WRITE_RATE_LIMIT", 5),
		Burst:             getEnvInt("PROJECTS_WRITE_BURST_LIMIT", 10),
		TTL:               getEnvDuration("PROJECTS_RATE_LIMIT_TTL", 10*time.Minute),
		CleanupInterval:   getEnvDuration("PROJECTS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})
	
	readLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("PROJECTS_READ_RATE_LIMIT", 15),
		Burst:             getEnvInt("PROJECTS_READ_BURST_LIMIT", 30),
		TTL:               getEnvDuration("PROJECTS_RATE_LIMIT_TTL", 10*time.Minute),
		CleanupInterval:   getEnvDuration("PROJECTS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})
	
	queryLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("PROJECTS_QUERY_RATE_LIMIT", 8),
		Burst:             getEnvInt("PROJECTS_QUERY_BURST_LIMIT", 15),
		TTL:               getEnvDuration("PROJECTS_RATE_LIMIT_TTL", 10*time.Minute),
		CleanupInterval:   getEnvDuration("PROJECTS_RATE_LIMIT_CLEANUP", 5*time.Minute),
	})

	writeRoutes := r.Group("/projects")
	writeRoutes.Use(writeLimiter.RateLimitMiddleware())
	{
		writeRoutes.POST("", createProjectController.Execute)
		writeRoutes.PUT("/:id", updateProjectController.Execute)
		writeRoutes.DELETE("/:id", deleteProjectController.Execute)
	}

	readRoutes := r.Group("/projects")
	readRoutes.Use(readLimiter.RateLimitMiddleware())
	{
		readRoutes.GET("", getProjectsController.Execute)
		readRoutes.GET("/id/:id", getProjectByIdController.Execute)
		readRoutes.GET("/user/:userId", getProjectByUserId.Execute)
	}

	queryRoutes := r.Group("/projects")
	queryRoutes.Use(queryLimiter.RateLimitMiddleware())
	{
		queryRoutes.GET("/nombre/:nombre", getProjectByNameController.Execute)
		queryRoutes.GET("/categoria/:categoria", getProjectByCategoryController.Execute)
		queryRoutes.GET("/fecha/:fecha", getProjectByDateController.Execute)
		queryRoutes.GET("/stats", getProjectsStats.Execute)
	}
}