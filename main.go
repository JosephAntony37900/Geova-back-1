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
	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: engine,
	}

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("🚀 Servidor iniciando en %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Canal para señales de sistema operativo
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Esperar señal de shutdown o error del servidor
	select {
	case err := <-serverErrors:
		log.Fatalf("Error al iniciar el servidor: %v", err)
	case sig := <-shutdown:
		log.Printf("Señal de shutdown recibida: %v", sig)

		// Crear contexto con timeout para el shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Apagar servidor HTTP
		log.Println("Apagando servidor HTTP...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error durante shutdown del servidor: %v", err)
			srv.Close()
		}

		// Cerrar infraestructura de proyectos (incluye ImageUploadWorkerService)
		projectInfra.Shutdown()

		log.Println("Servidor apagado correctamente")
	}
}
