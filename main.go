package main

import (
	"log"

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
	project_infra.InitProjectDependencies(engine)

	// Iniciar servidor
	port := "0.0.0.0:8000"
	log.Printf("🚀 Servidor iniciando en %s", port)
	
	if err := engine.Run(port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}