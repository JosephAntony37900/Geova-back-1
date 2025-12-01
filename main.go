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

	engine := gin.Default()

	engine.Use(core.SetupCORS())

	user_infra.InitUserDependencies(engine)
	projectInfra := project_infra.InitProjectDependencies(engine)

	port := "0.0.0.0:8000"
	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	go startServer(srv, port)

	waitForShutdown(srv, projectInfra)
}

// startServer inicia el servidor HTTP en una goroutine
func startServer(srv *http.Server, port string) {
	log.Printf("🚀 Servidor iniciando en %s", port)
	
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Error al iniciar el servidor: %v", err)
	}
}

func waitForShutdown(srv *http.Server, projectInfra *project_infra.ProjectInfrastructure) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit
	log.Println("Señal de shutdown recibida, cerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Shutdown de infraestructura de proyectos 
	log.Println("Cerrando infraestructura de proyectos...")
	projectInfra.Shutdown()

	// Shutdown del servidor HTTP
	log.Println("Cerrando servidor HTTP...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Error durante shutdown del servidor: %v", err)
	}

	log.Println("Servidor cerrado exitosamente")
}