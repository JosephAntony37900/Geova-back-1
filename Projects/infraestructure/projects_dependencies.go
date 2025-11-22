// Geova-back-1/Projects/infraestructure/projects_dependencies.go
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

type ProjectInfrastructure struct {
	DB                 *core.Conn_MySQL
	ProjectRepo        domain_projects.ProjectRepository
	ImageUploadService *domain_services.ImageUploadWorkerService // ✅ Worker service
}

func NewProjectInfrastructure() *ProjectInfrastructure {
	db := core.NewDatabaseConnection()

	if db == nil || db.DB == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar la conexión a la base de datos")
	}

	log.Println("INFO: Conexión a base de datos establecida")

	projectRepo := repo_projects.NewProjectMySQLRepository(db)

	return &ProjectInfrastructure{
		DB:          db,
		ProjectRepo: projectRepo,
	}
}

func InitProjectDependencies(engine *gin.Engine) *ProjectInfrastructure {
	log.Println("INFO: Inicializando infraestructura de proyectos...")

	// Crear infraestructura base
	infrastructure := NewProjectInfrastructure()

	// Inicializar Cloudinary
	log.Println("INFO: Inicializando servicio de Cloudinary...")
	cloudinaryAdapter, err := services_projects.NewCloudinaryAdapter()
	if err != nil {
		log.Printf("ERROR: No se pudo inicializar Cloudinary: %v", err)
		panic("Error crítico al inicializar Cloudinary: " + err.Error())
	}
	log.Println("Cloudinary inicializado exitosamente")

	log.Println("INFO: Inicializando Worker Service de imágenes...")
	imageUploadService := domain_services.NewImageUploadWorkerService(
		cloudinaryAdapter,
		infrastructure.ProjectRepo,
		3, 
	)
	infrastructure.ImageUploadService = imageUploadService
	log.Println("Worker Service inicializado con 3 workers")

	log.Println("INFO: Inicializando casos de uso...")
	createProjectUseCase := app_projects.NewCreateProjectUseCase(
		infrastructure.ProjectRepo,
		imageUploadService, // 
	)
	getAllProjectsUseCase := app_projects.NewGeProjectsUseCase(infrastructure.ProjectRepo)
	getProjectByIdUseCase := app_projects.NewGetProjectByIdUseCase(infrastructure.ProjectRepo)
	getProjectByNameUseCase := app_projects.NewGetProjectsByNameUseCase(infrastructure.ProjectRepo)
	getProjectByCategoryUseCase := app_projects.NewGetProjectsByCategoryUseCase(infrastructure.ProjectRepo)
	getProjectByDateUseCase := app_projects.NewGetProjectsByDateUseCase(infrastructure.ProjectRepo)
	updateProjectUseCase := app_projects.NewUpdateProjectUseCase(
		infrastructure.ProjectRepo,
		imageUploadService, //Inyectar worker service
	)
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

	return infrastructure
}

// Shutdown cierra todas las conexiones de forma limpia
func (pi *ProjectInfrastructure) Shutdown() {

	if pi.ImageUploadService != nil {
		pi.ImageUploadService.Shutdown()
	}

	if pi.DB != nil && pi.DB.DB != nil {
		pi.DB.DB.Close()
		log.Println("INFO: Conexión a base de datos cerrada")
	}

}