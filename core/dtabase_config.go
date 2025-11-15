// core/database_config.go
package core

import (
	"fmt"
	"log"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Schema   string
	Port     string
}

// Configuración para BD remota
func GetDatabaseConfig() DatabaseConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error al cargar .env: %v", err)
	}

	return DatabaseConfig{
		Host:     os.Getenv("REMOTE_DB_HOST"),
		User:     os.Getenv("REMOTE_DB_USER"),
		Password: os.Getenv("REMOTE_DB_PASS"),
		Schema:   os.Getenv("REMOTE_DB_SCHEMA"),
		Port:     os.Getenv("REMOTE_DB_PORT"),
	}
}

// Crear conexión a base de datos
func CreateDBConnection(config DatabaseConfig) *Conn_MySQL {
	port := config.Port
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		config.User, config.Password, config.Host, port, config.Schema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error al abrir la base de datos: %v", err)
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al abrir la base de datos: %v", err)}
	}

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Probar la conexión
	if err := db.Ping(); err != nil {
		log.Printf("Error al verificar la conexión: %v", err)
		db.Close()
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al verificar la conexión: %v", err)}
	}

	log.Printf("Conexión exitosa a base de datos")
	return &Conn_MySQL{DB: db, Err: ""}
}

// Inicializar conexión a base de datos
func NewDatabaseConnection() *Conn_MySQL {
	config := GetDatabaseConfig()
	db := CreateDBConnection(config)
	
	if db.DB == nil {
		log.Fatal("ERROR CRÍTICO: No se puede conectar a la base de datos")
	}
	
	return db
}