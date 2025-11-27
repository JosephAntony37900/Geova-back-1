// geova-back-1/Users/infraestructure/routes/users_routers.go
package routes

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/controllers"
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

// RateLimitMiddleware crea el middleware de rate limiting
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



func SetupUserRoutes(r *gin.Engine, 
	createUserController *controllers.CreateUserController,
	getUsersController *controllers.GetAllUsersController,
	getUsersControllerById *controllers.GetUserByIdController,
	updateUserController *controllers.UpdateUserController,
	deleteUserController *controllers.DeleteUserController,
	loginUserController *controllers.LoginUserController,
) {
	
	// 4 intentos/seg, burst 3 (solo 3 intentos rápidos)
	loginLimiter := NewRateLimiter(4, 3)
	
	// 1 registro cada 2 segundos, burst 2
	registerLimiter := NewRateLimiter(0.5, 2)
	
	// 3 peticiones/seg, burst 5
	modifyLimiter := NewRateLimiter(3, 5)
	
	// 20 peticiones/seg, burst 40
	readLimiter := NewRateLimiter(20, 40)

	// Login
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

	// Modificación
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