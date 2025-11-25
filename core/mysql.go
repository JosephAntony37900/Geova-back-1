//geova-back-1/core/mysql.go
package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Conn_MySQL struct {
	DB  *sql.DB
	Err string
}

// Propuesta de valores optimizados para el pool de conexiones MySQL
func ConfigureDBPool(db *sql.DB) {
	// Máximo de conexiones abiertas simultáneas
	
	db.SetMaxOpenConns(100)
	
	// Conexiones inactivas mantenidas en el pool
	db.SetMaxIdleConns(25)
	
	// Tiempo máximo de vida de una conexión
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Tiempo máximo que una conexión puede estar inactiva
	db.SetConnMaxIdleTime(10 * time.Minute)
	
	log.Println("INFO: Pool de conexiones configurado - MaxOpen:100, MaxIdle:25, MaxLifetime:5m, MaxIdleTime:10m")
}

func GetDBPool() *Conn_MySQL {
	error := ""
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	dbHost := os.Getenv("REMOTE_DB_HOST")
	dbUser := os.Getenv("REMOTE_DB_USER")
	dbPass := os.Getenv("REMOTE_DB_PASS")
	dbSchema := os.Getenv("REMOTE_DB_SCHEMA")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPass, dbHost, dbSchema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		error = fmt.Sprintf("error al abrir la base de datos: %v", err)
		return &Conn_MySQL{DB: nil, Err: error}
	}

	ConfigureDBPool(db)

	if err := db.Ping(); err != nil {
		db.Close()
		error = fmt.Sprintf("error al verificar la conexión a la base de datos: %v", err)
		return &Conn_MySQL{DB: nil, Err: error}
	}

	return &Conn_MySQL{DB: db, Err: error}
}

func (conn *Conn_MySQL) ExecutePreparedQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := conn.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta preparada: %w", err)
	}

	return result, nil
}

func (conn *Conn_MySQL) FetchRows(query string, values ...interface{}) *sql.Rows {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		fmt.Printf("error al ejecutar la consulta SELECT: %v", err)
	}

	return rows
}