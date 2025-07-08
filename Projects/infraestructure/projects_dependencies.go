package infraestructure


import (
	_"database/sql"
	app_projects "github.com/JosephAntony37900/Geova-back-1/Projects/application"
	control_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
	repo_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/repository"
	routes_projects "github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/routes"
	"github.com/JosephAntony37900/Geova-back-1/core"
	"github.com/gin-gonic/gin"
)

func InitprojectDependencies(engine *gin.Engine, conn *core.Conn_MySQL) {

	projectRepo := repo_projects.NewProjectMySQLRepository(conn)
	
	
	createProjectUseCase := app_projects.NewCreateProjectUseCase(projectRepo )
	getAllProjectsUseCase := app_projects.NewGeProjectsUseCase(projectRepo)
	getProjectByIdUseCase := app_projects.NewGetProjectByIdUseCase(projectRepo)
	getProjectByNameUseCase := app_projects.NewGetProjectsByNameUseCase(projectRepo)
	getProjectByCategoryUseCase := app_projects.NewGetProjectsByCategoryUseCase(projectRepo)
	getProjectByDateUseCase := app_projects.NewGetProjectsByDateUseCase(projectRepo)
	upateProjectUseCase := app_projects.NewUpdateProjectUseCase(projectRepo )
	deleteProjectUseCase := app_projects.NewDeleteProjectUseCase(projectRepo)


	createProjectController := control_projects.NewCreateProjectController(createProjectUseCase)
	getAllProjectController := control_projects.NewGetAllProjectsController(getAllProjectsUseCase)
	getByIdProjectController := control_projects.NewGetProjectByIdUseController(getProjectByIdUseCase)
	getProjectByNameController := control_projects.NewGetProjectByNameController(getProjectByNameUseCase)
	getProjectByCategoryController := control_projects.NewGetProjectByCategoryController(getProjectByCategoryUseCase)
	getProjectByDateController := control_projects.NewGetProjectByDateController(getProjectByDateUseCase)
	updateProjectController := control_projects.NewUpdateProjectController(upateProjectUseCase)
	deleteProjectController := control_projects.NewDeleteProjectController(deleteProjectUseCase)


	

	routes_projects.SetUpProjectsRoutes(engine, createProjectController, getAllProjectController, getByIdProjectController, getProjectByNameController, getProjectByCategoryController, getProjectByDateController,updateProjectController, deleteProjectController )
}