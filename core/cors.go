package core

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Lista de orígenes permitidos
			allowedOrigins := []string{
				"https://www.geova.pro",
				"https://geova.pro",
				"http://localhost:3000",
				"http://localhost:5173",
				"http://localhost:4200",
			}
			
			// Verificar si hay un origen personalizado en las variables de entorno
			customOrigin := os.Getenv("ALLOWED_ORIGIN")
			if customOrigin != "" {
				allowedOrigins = append(allowedOrigins, customOrigin)
			}
			
			// Comprobar si el origen está en la lista permitida
			for _, allowed := range allowedOrigins {
				if strings.EqualFold(origin, allowed) {
					return true
				}
			}
			
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
