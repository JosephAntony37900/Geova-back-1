package infraestructure

import (
	"log"

	app_projects "github.com/JosephAntony37900/Geova-back-1/Projects/application"
	control_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
	domain_projects "github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	repo_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/repository"
	routes_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/routes"
	services_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/services/adapters"
	"github.com/JosephAntony37900/Geova-back-1/core"

	"github.com/gin-gonic/gin"
)

// ProjectInfrastructure encapsula toda la infraestructura de proyectos
type ProjectInfrastructure struct {
	DatabaseManager *core.DatabaseManager
	ProjectRepo     domain_projects.ProjectRepository
}

// NewProjectInfrastructure crea e inicializa toda la infraestructura de proyectos
func NewProjectInfrastructure() *ProjectInfrastructure {
	// Inicializar el DatabaseManager (maneja ambas conexiones)
	dbManager := core.NewDatabaseManager()
	
	// Validar que el DatabaseManager se inicializó correctamente
	if dbManager == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar el DatabaseManager")
	}
	
	// Verificar estado de las conexiones
	if dbManager.LocalDB == nil || dbManager.LocalDB.DB == nil {
		panic("ERROR CRÍTICO: No se puede inicializar sin conexión a BD local")
	}
	
	// Log del estado de las conexiones
	if dbManager.RemoteDB == nil || dbManager.RemoteDB.DB == nil {
		log.Println("INFO: Iniciando en modo offline - solo BD local disponible")
		log.Println("INFO: Los datos se sincronizarán automáticamente cuando la BD remota esté disponible")
	} else {
		log.Println("INFO: Iniciando con ambas conexiones disponibles (local y remota)")
	}
	
	// Crear repositorio usando el DatabaseManager
	projectRepo := repo_projects.NewProjectMySQLRepository(
		dbManager.LocalDB, 
		dbManager.RemoteDB,
	)
	
	return &ProjectInfrastructure{
		DatabaseManager: dbManager,
		ProjectRepo:     projectRepo,
	}
}

// InitProjectDependencies inicializa todas las dependencias y configura las rutas
func InitProjectDependencies(engine *gin.Engine) *ProjectInfrastructure {
	log.Println("INFO: Inicializando infraestructura de proyectos...")
	
	// Crear infraestructura
	infrastructure := NewProjectInfrastructure()
	
	// Inicializar Cloudinary
	log.Println("INFO: Inicializando servicio de Cloudinary...")
	cloudinaryAdapter, err := services_projects.NewCloudinaryAdapter()
	if err != nil {
		log.Printf("ERROR: No se pudo inicializar Cloudinary: %v", err)
		panic("Error crítico al inicializar Cloudinary: " + err.Error())
	}
	log.Println("INFO: Cloudinary inicializado exitosamente")

	// Crear casos de uso
	log.Println("INFO: Inicializando casos de uso...")
	createProjectUseCase := app_projects.NewCreateProjectUseCase(infrastructure.ProjectRepo, cloudinaryAdapter)
	getAllProjectsUseCase := app_projects.NewGeProjectsUseCase(infrastructure.ProjectRepo)
	getProjectByIdUseCase := app_projects.NewGetProjectByIdUseCase(infrastructure.ProjectRepo)
	getProjectByNameUseCase := app_projects.NewGetProjectsByNameUseCase(infrastructure.ProjectRepo)
	getProjectByCategoryUseCase := app_projects.NewGetProjectsByCategoryUseCase(infrastructure.ProjectRepo)
	getProjectByDateUseCase := app_projects.NewGetProjectsByDateUseCase(infrastructure.ProjectRepo)
	updateProjectUseCase := app_projects.NewUpdateProjectUseCase(infrastructure.ProjectRepo, cloudinaryAdapter)
	deleteProjectUseCase := app_projects.NewDeleteProjectUseCase(infrastructure.ProjectRepo)
	getProjectsByUserIdUseCase := app_projects.NewGetProjectsByUserIdUseCase(infrastructure.ProjectRepo)

	// Crear controladores
	log.Println("INFO: Inicializando controladores...")
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
	log.Println("INFO: Configurando rutas de proyectos...")
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
	
	log.Println("INFO: Infraestructura de proyectos inicializada exitosamente")
	return infrastructure
}

// Shutdown cierra todas las conexiones de forma limpia
func (pi *ProjectInfrastructure) Shutdown() {
	log.Println("INFO: Cerrando infraestructura de proyectos...")
	
	if pi.DatabaseManager != nil {
		pi.DatabaseManager.Close()
	}
	
	log.Println("INFO: Infraestructura de proyectos cerrada exitosamente")
}

// GetConnectionStatus retorna el estado de las conexiones
func (pi *ProjectInfrastructure) GetConnectionStatus() map[string]bool {
	status := make(map[string]bool)
	
	// Verificar conexión local
	status["local"] = false
	if pi.DatabaseManager.LocalDB != nil && pi.DatabaseManager.LocalDB.DB != nil {
		if err := pi.DatabaseManager.LocalDB.DB.Ping(); err == nil {
			status["local"] = true
		}
	}
	
	// Verificar conexión remota
	status["remote"] = false
	if pi.DatabaseManager.RemoteDB != nil && pi.DatabaseManager.RemoteDB.DB != nil {
		if err := pi.DatabaseManager.RemoteDB.DB.Ping(); err == nil {
			status["remote"] = true
		}
	}
	
	return status
}

// ReconnectRemoteDB intenta reconectar a la BD remota
func (pi *ProjectInfrastructure) ReconnectRemoteDB() bool {
	if pi.DatabaseManager != nil {
		pi.DatabaseManager.ReconnectRemote()
		
		// Verificar si la reconexión fue exitosa
		if pi.DatabaseManager.RemoteDB != nil && pi.DatabaseManager.RemoteDB.DB != nil {
			if err := pi.DatabaseManager.RemoteDB.DB.Ping(); err == nil {
				log.Println("INFO: Reconexión a BD remota exitosa")
				return true
			}
		}
	}
	
	log.Println("WARNING: No se pudo reconectar a la BD remota")
	return false
}

// HealthCheck verifica el estado general de la infraestructura
func (pi *ProjectInfrastructure) HealthCheck() map[string]interface{} {
	healthStatus := make(map[string]interface{})
	
	// Estado de conexiones
	connectionStatus := pi.GetConnectionStatus()
	healthStatus["connections"] = connectionStatus
	
	// Estado general
	healthStatus["healthy"] = connectionStatus["local"] // Mínimo requerido es la BD local
	healthStatus["mode"] = "offline"
	if connectionStatus["remote"] {
		healthStatus["mode"] = "online"
	}
	
	// Información adicional
	healthStatus["sync_enabled"] = connectionStatus["remote"]
	
	return healthStatus
}