//geova-back-1/Projects/infraestructure/projects_dependencies.go
package infraestructure

import (
	"log"

	app_projects "github.com/JosephAntony37900/Geova-back-1/Projects/application"
	domain_projects "github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	domain_services "github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
	control_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
	repo_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/repository"
	routes_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/routes"
	services_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/services/adapters"
	"github.com/JosephAntony37900/Geova-back-1/core"

	"github.com/gin-gonic/gin"
)

// ProjectInfrastructure encapsula toda la infraestructura de proyectos
type ProjectInfrastructure struct {
	DB          *core.Conn_MySQL
	ProjectRepo domain_projects.ProjectRepository
	WorkerSrv   *domain_services.ImageUploadWorkerService
}

// NewProjectInfrastructure crea e inicializa toda la infraestructura de proyectos
func NewProjectInfrastructure() *ProjectInfrastructure {
	// Inicializar conexión a base de datos
	db := core.NewDatabaseConnection()

	if db == nil || db.DB == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar la conexión a la base de datos")
	}

	log.Println("INFO: Conexión a base de datos establecida")

	// Crear repositorio
	projectRepo := repo_projects.NewProjectMySQLRepository(db)

	return &ProjectInfrastructure{
		DB:          db,
		ProjectRepo: projectRepo,
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

	// Inicializar ImageUploadWorkerService
	log.Println("INFO: Inicializando ImageUploadWorkerService...")
	workerService := domain_services.NewImageUploadWorkerService(cloudinaryAdapter, 3, 100)
	infrastructure.WorkerSrv = workerService
	log.Println("INFO: ImageUploadWorkerService inicializado exitosamente")

	// Crear casos de uso
	log.Println("INFO: Inicializando casos de uso...")
	createProjectUseCase := app_projects.NewCreateProjectUseCase(infrastructure.ProjectRepo, cloudinaryAdapter, workerService)
	getAllProjectsUseCase := app_projects.NewGeProjectsUseCase(infrastructure.ProjectRepo)
	getProjectByIdUseCase := app_projects.NewGetProjectByIdUseCase(infrastructure.ProjectRepo)
	getProjectByNameUseCase := app_projects.NewGetProjectsByNameUseCase(infrastructure.ProjectRepo)
	getProjectByCategoryUseCase := app_projects.NewGetProjectsByCategoryUseCase(infrastructure.ProjectRepo)
	getProjectByDateUseCase := app_projects.NewGetProjectsByDateUseCase(infrastructure.ProjectRepo)
	getProjectStatsUseCase := app_projects.NewGetProjectStatsUseCase(infrastructure.ProjectRepo)
	updateProjectUseCase := app_projects.NewUpdateProjectUseCase(infrastructure.ProjectRepo, cloudinaryAdapter, workerService)
	deleteProjectUseCase := app_projects.NewDeleteProjectUseCase(infrastructure.ProjectRepo)
	getProjectsByUserIdUseCase := app_projects.NewGetProjectsByUserIdUseCase(infrastructure.ProjectRepo)
	getTotalProjectsByUserUseCase := app_projects.NewGetTotalProjectsByUserUseCase(infrastructure.ProjectRepo)


	// Crear controladores
	log.Println("INFO: Inicializando controladores...")
	createProjectController := control_projects.NewCreateProjectController(createProjectUseCase)
	getAllProjectController := control_projects.NewGetAllProjectsController(getAllProjectsUseCase)
	getByIdProjectController := control_projects.NewGetProjectByIdUseController(getProjectByIdUseCase)
	getProjectByNameController := control_projects.NewGetProjectByNameController(getProjectByNameUseCase)
	getProjectByCategoryController := control_projects.NewGetProjectByCategoryController(getProjectByCategoryUseCase)
	getProjectByDateController := control_projects.NewGetProjectByDateController(getProjectByDateUseCase)
	getProjectsStatsController := control_projects.NewGetProjectStatsController(getProjectStatsUseCase)
	updateProjectController := control_projects.NewUpdateProjectController(updateProjectUseCase)
	deleteProjectController := control_projects.NewDeleteProjectController(deleteProjectUseCase)
	getProjectsByUserIdController := control_projects.NewGetProjectsByUserIdController(getProjectsByUserIdUseCase)
	getTotalProjectsByUserController := control_projects.NewGetTotalProjectsByUserController(getTotalProjectsByUserUseCase)

	// Configurar rutas
	log.Println("INFO: Configurando rutas de proyectos...")
	routes_projects.SetUpProjectsRoutes(engine,
		createProjectController,
		getAllProjectController,
		getByIdProjectController,
		getProjectByNameController,
		getProjectByCategoryController,
		getProjectByDateController,
		getProjectsStatsController,
		updateProjectController,
		deleteProjectController,
		getProjectsByUserIdController,
		getTotalProjectsByUserController,)

	log.Println("INFO: Infraestructura de proyectos inicializada exitosamente")
	return infrastructure
}

// Shutdown cierra todas las conexiones de forma limpia
func (pi *ProjectInfrastructure) Shutdown() {
	log.Println("INFO: Cerrando infraestructura de proyectos...")

	// Shutdown del worker service primero
	if pi.WorkerSrv != nil {
		pi.WorkerSrv.Shutdown()
		log.Println("INFO: Worker service cerrado")
	}

	if pi.DB != nil && pi.DB.DB != nil {
		pi.DB.DB.Close()
		log.Println("INFO: Conexión a base de datos cerrada")
	}

	log.Println("INFO: Infraestructura de proyectos cerrada exitosamente")
}
