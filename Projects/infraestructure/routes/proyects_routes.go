package routes

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
)

// RateLimiter gestiona los limitadores por IP
type RateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter obtiene o crea un limiter para una IP
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.ips[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)
		
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Demasiadas peticiones. Intenta más tarde.",
				"message": "Rate limit exceeded",
				"ip":      ip,
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

func SetUpProjectsRoutes(r *gin.Engine, 
	createProjectController *controllers.CreateProjectController,
	getProjectsController *controllers.GetAllProjectsController,
	getProjectByIdController *controllers.GetProjectByIdController,
	getProjectByNameController *controllers.GetProjectByNameController,
	GetProjectByCategoryController *controllers.GetProjectByCategoryController,
	getProjetcByDateController *controllers.GetProjectByDateController,
	getProjectsStats *controllers.GetProjectStatsController,
	updateProjectController *controllers.UpdateProjectController,
	deleteProjectController *controllers.DeleteProjectController,
	getProjectByUserId *controllers.GetProjectsByUserIdController,
) {
	// Operaciones de escritura: moderado (5 req/seg, burst de 10)
	writeLimiter := NewRateLimiter(5, 10)
	
	// Operaciones de lectura: más permisivo (15 req/seg, burst de 30)
	readLimiter := NewRateLimiter(15, 30)
	
	// Stats y consultas complejas: intermedio (8 req/seg, burst de 15)
	queryLimiter := NewRateLimiter(8, 15)

	// Rutas de escritura  - Más restrictivas
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
		queryRoutes.GET("/categoria/:categoria", GetProjectByCategoryController.Execute)
		queryRoutes.GET("/fecha/:fecha", getProjetcByDateController.Execute)
		queryRoutes.GET("/stats", getProjectsStats.Execute)
	}
}

