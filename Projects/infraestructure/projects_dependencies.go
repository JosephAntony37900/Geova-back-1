package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	app_projects "github.com/JosephAntony37900/Geova-back-1/Projects/application"
	control_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
	repo_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/repository"
	routes_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/routes"
	services_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/services/adapters"
	"github.com/JosephAntony37900/Geova-back-1/core"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Función para crear conexión a BD local
func createLocalDBConnection() *core.Conn_MySQL {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No se pudo cargar archivo .env para BD local: %v", err)
	}

	// Variables para BD local (SQLite o MySQL local)
	localDBHost := os.Getenv("LOCAL_DB_HOST")
	localDBUser := os.Getenv("LOCAL_DB_USER") 
	localDBPass := os.Getenv("LOCAL_DB_PASS")
	localDBSchema := os.Getenv("LOCAL_DB_SCHEMA")

	// Si no están definidas, usar valores por defecto para BD local
	if localDBHost == "" {
		localDBHost = "localhost"
	}
	if localDBUser == "" {
		localDBUser = "root"
	}
	if localDBSchema == "" {
		localDBSchema = "geova_local"
	}

	// Crear conexión local
	localConn := &core.Conn_MySQL{}
	
	// Puedes usar la función GetDBPool existente o crear una específica para local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", localDBUser, localDBPass, localDBHost, localDBSchema)
	db, err := sql.Open("mysql", dsn)
	
	if err != nil {
		log.Printf("Error al conectar a BD local: %v", err)
		return nil
	}

	db.SetMaxOpenConns(5) // Menos conexiones para BD local
	
	if err := db.Ping(); err != nil {
		log.Printf("Warning: BD local no disponible: %v", err)
		db.Close()
		return nil
	}

	localConn.DB = db
	log.Println("Conexión a BD local establecida exitosamente")
	return localConn
}

// Función para crear conexión a BD remota
func createRemoteDBConnection() *core.Conn_MySQL {
	remoteConn := core.GetDBPool() // Tu función existente para BD remota
	
	if remoteConn.Err != "" {
		log.Printf("Warning: BD remota no disponible: %s", remoteConn.Err)
		return nil
	}
	
	log.Println("Conexión a BD remota establecida exitosamente")
	return remoteConn
}

func InitProjectDependencies(engine *gin.Engine) {
	// Crear conexiones a ambas bases de datos
	localConn := createLocalDBConnection()
	remoteConn := createRemoteDBConnection()
	
	// Validar que al menos la BD local esté disponible
	if localConn == nil {
		panic("ERROR CRÍTICO: No se puede inicializar sin BD local")
	}
	
	// Si la BD remota no está disponible, continuar solo con local
	if remoteConn == nil {
		log.Println("INFO: Iniciando en modo offline - solo BD local disponible")
	}

	// Crear repositorio con ambas conexiones
	projectRepo := repo_projects.NewProjectMySQLRepository(localConn, remoteConn)
	
	// Inicializar Cloudinary
	cloudinaryAdapter, err := services_projects.NewCloudinaryAdapter()
	if err != nil {
		panic("Error al inicializar Cloudinary: " + err.Error())
	}

	// Crear casos de uso
	createProjectUseCase := app_projects.NewCreateProjectUseCase(projectRepo, cloudinaryAdapter)
	getAllProjectsUseCase := app_projects.NewGeProjectsUseCase(projectRepo)
	getProjectByIdUseCase := app_projects.NewGetProjectByIdUseCase(projectRepo)
	getProjectByNameUseCase := app_projects.NewGetProjectsByNameUseCase(projectRepo)
	getProjectByCategoryUseCase := app_projects.NewGetProjectsByCategoryUseCase(projectRepo)
	getProjectByDateUseCase := app_projects.NewGetProjectsByDateUseCase(projectRepo)
	updateProjectUseCase := app_projects.NewUpdateProjectUseCase(projectRepo, cloudinaryAdapter)
	deleteProjectUseCase := app_projects.NewDeleteProjectUseCase(projectRepo)
	getProjectsByUserIdUseCase := app_projects.NewGetProjectsByUserIdUseCase(projectRepo)

	// Crear controladores
	createProjectController := control_projects.NewCreateProjectController(createProjectUseCase)
	getAllProjectController := control_projects.NewGetAllProjectsController(getAllProjectsUseCase)
	getByIdProjectController := control_projects.NewGetProjectByIdUseController(getProjectByIdUseCase)
	getProjectByNameController := control_projects.NewGetProjectByNameController(getProjectByNameUseCase)
	getProjectByCategoryController := control_projects.NewGetProjectByCategoryController(getProjectByCategoryUseCase)
	getProjectByDateController := control_projects.NewGetProjectByDateController(getProjectByDateUseCase)
	updateProjectController := control_projects.NewUpdateProjectController(updateProjectUseCase)
	deleteProjectController := control_projects.NewDeleteProjectController(deleteProjectUseCase)
	getProjectsByUserIdController := control_projects.NewGetProjectsByUserIdController(getProjectsByUserIdUseCase)

	// Configurar rutas
	routes_projects.SetUpProjectsRoutes(engine, 
		createProjectController, 
		getAllProjectController, 
		getByIdProjectController, 
		getProjectByNameController, 
		getProjectByCategoryController, 
		getProjectByDateController,
		updateProjectController, 
		deleteProjectController,
		getProjectsByUserIdController)
}