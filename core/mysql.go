package core;

import (
	"fmt"
	"log"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

type Conn_MySQL struct {
	DB *sql.DB
	Err string
}


func GetDBPool() *Conn_MySQL {

	error := ""
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	// Obtener las variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSchema := os.Getenv("DB_SCHEMA")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbSchema)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		error = fmt.Sprintf("error al abrir la base de datos: %w", err)
	}

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)

	// Probar la conexión
	if err := db.Ping(); err != nil {
		db.Close()
		error = fmt.Sprintf("error al verificar la conexión a la base de datos: %w", err)
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


func (conn *Conn_MySQL) FetchRows(query string, values ...interface{}) (*sql.Rows) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		fmt.Printf ("error al ejecutar la consulta SELECT: %w", err)
	}

	return rows
}