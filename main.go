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

	//rabbitmqUser := os.Getenv("RABBITMQ_USER")
	//rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	//rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	//rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	//rabbitmqURI := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)

	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error inicializando la conexión a MySQL: %v", conn.Err)
	}
	defer conn.DB.Close()

	engine := gin.Default()
	engine.Use(core.SetupCORS())

	user_infra.InitUserDependencies(engine, conn)
	project_infra.InitProjectDependencies(engine)

	engine.Run("0.0.0.0:8000")
}