package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/core"
	user_infra "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure"
	project_infra "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure"

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

	// Configurar servidor
	port := "0.0.0.0:8000"
	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	// Goroutine para manejar señales de shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		sig := <-sigChan
		log.Printf("Señal recibida: %v. Iniciando shutdown graceful...", sig)

		// Shutdown worker service and other infrastructure
		projectInfra.Shutdown()

		// Shutdown HTTP server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("ERROR: Forzando cierre del servidor: %v", err)
		}

		log.Println("Shutdown completado exitosamente")
	}()

	// Iniciar servidor
	log.Printf("🚀 Servidor iniciando en %s", port)
	
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}