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

type DatabaseManager struct {
	LocalDB  *Conn_MySQL
	RemoteDB *Conn_MySQL
}

// Configuración para BD local (Raspberry Pi)
func GetLocalDatabaseConfig() DatabaseConfig {
		err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error al cargar .env para BD local: %v", err)
	}
	return DatabaseConfig{
		
		Host:     os.Getenv("LOCAL_DB_HOST"),
		User:     os.Getenv("LOCAL_DB_USER"),      // Cambiar según tu configuración
		Password: os.Getenv("LOCAL_DB_PASS"),  // Cambiar según tu configuración
		Schema:   os.Getenv("LOCAL_DB_SCHEMA"),     // Cambiar según tu configuración
		Port:     os.Getenv("PORT"),
	}
}

// Configuración para BD remota (desplegada)
func GetRemoteDatabaseConfig() DatabaseConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error al cargar .env para BD remota: %v", err)
	}

	return DatabaseConfig{
		Host:     os.Getenv("REMOTE_DB_HOST"),
		User:     os.Getenv("REMOTE_DB_USER"),
		Password: os.Getenv("REMOTE_DB_PASS"),
		Schema:   os.Getenv("REMOTE_DB_SCHEMA"),
		Port:     os.Getenv("REMOTE_DB_PORT"),
	}
}

// Crear conexión a base de datos con configuración específica
func CreateDBConnection(config DatabaseConfig) *Conn_MySQL {
	port := config.Port
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		config.User, config.Password, config.Host, port, config.Schema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error al abrir la base de datos %s: %v", config.Host, err)
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al abrir la base de datos: %v", err)}
	}

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Probar la conexión
	if err := db.Ping(); err != nil {
		log.Printf("Error al verificar la conexión a %s: %v", config.Host, err)
		db.Close()
		return &Conn_MySQL{DB: nil, Err: fmt.Sprintf("error al verificar la conexión: %v", err)}
	}

	log.Printf("Conexión exitosa a base de datos: %s", config.Host)
	return &Conn_MySQL{DB: db, Err: ""}
}

// Inicializar manager de base de datos con conexiones local y remota
func NewDatabaseManager() *DatabaseManager {
	// Crear conexión local
	localConfig := GetLocalDatabaseConfig()
	localDB := CreateDBConnection(localConfig)
	
	// Crear conexión remota
	remoteConfig := GetRemoteDatabaseConfig()
	remoteDB := CreateDBConnection(remoteConfig)
	
	// Verificar que al menos la BD local esté disponible
	if localDB.DB == nil {
		log.Fatal("ERROR CRÍTICO: No se puede conectar a la base de datos local")
	}
	
	if remoteDB.DB == nil {
		log.Printf("WARNING: No se puede conectar a la base de datos remota. Funcionando en modo offline.")
		remoteDB = nil
	}
	
	return &DatabaseManager{
		LocalDB:  localDB,
		RemoteDB: remoteDB,
	}
}

// Método para reconectar a la BD remota si está disponible
func (dm *DatabaseManager) ReconnectRemote() {
	if dm.RemoteDB != nil && dm.RemoteDB.DB != nil {
		return // Ya está conectada
	}
	
	remoteConfig := GetRemoteDatabaseConfig()
	remoteDB := CreateDBConnection(remoteConfig)
	
	if remoteDB.DB != nil {
		dm.RemoteDB = remoteDB
		log.Println("INFO: Reconexión exitosa a base de datos remota")
	}
}

// Método para cerrar conexiones
func (dm *DatabaseManager) Close() {
	if dm.LocalDB != nil && dm.LocalDB.DB != nil {
		dm.LocalDB.DB.Close()
		log.Println("INFO: Conexión local cerrada")
	}
	
	if dm.RemoteDB != nil && dm.RemoteDB.DB != nil {
		dm.RemoteDB.DB.Close()
		log.Println("INFO: Conexión remota cerrada")
	}
}