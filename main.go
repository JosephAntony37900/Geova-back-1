package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	project_infra "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure"
	user_infra "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure"
	"github.com/JosephAntony37900/Geova-back-1/core"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error cargando el archivo .env: %v", err)
	}

	// Configurar Gin
	engine := gin.Default()

	// Configurar CORS
	engine.Use(core.SetupCORS())

	// Inicializar dependencias de usuarios y proyectos
	user_infra.InitUserDependencies(engine)
	projectInfra := project_infra.InitProjectDependencies(engine)

	// Configurar servidor HTTP
	port := "0.0.0.0:8000"
	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("🚀 Servidor iniciando en %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("INFO: Shutdown signal recibida, cerrando servidor...")

	// Shutdown graceful con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown de infraestructura de proyectos (incluye worker service)
	projectInfra.Shutdown()

	// Shutdown del servidor HTTP
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("ERROR: Error durante shutdown del servidor: %v", err)
	}

	log.Println("INFO: Servidor cerrado exitosamente")
}
