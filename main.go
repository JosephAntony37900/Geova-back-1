package main

import (
	_"os"
	"log"

	"github.com/JosephAntony37900/Geova-back-1/core"
	user_infra "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure"
	project_infra "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main () {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error inicializando la conexión a MySQL: %v", conn.Err)
	}
	defer conn.DB.Close()

	engine := gin.Default()
	engine.Use(core.SetupCORS())

	user_infra.InitUserDependencies(engine)
	project_infra.InitProjectDependencies(engine)

	engine.Run("0.0.0.0:8000")
}