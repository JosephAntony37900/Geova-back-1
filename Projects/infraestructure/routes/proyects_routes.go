package routes

import (
	"github.com/gin-gonic/gin"
	_"os"
	"github.com/JosephAntony37900/Geova-back-1/Projects/infraestructure/controllers"
)

func SetUpProjectsRoutes (r *gin.Engine, createProjectController *controllers.CreateProjectController, 
						getProjectsController *controllers.GetAllProjectsController,
						getProjectByIdController *controllers.GetProjectByIdController,
						updateProjectController *controllers.UpdateProjectController,
						deleteProjectController *controllers.DeleteProjectController) {

	r.POST("/projects", createProjectController.Execute)
	r.GET("/projects", getProjectsController.Execute)
	r.GET("/projects/:id", getProjectByIdController.Execute)
	r.PUT("/projects/:id", updateProjectController.Execute)
	r.DELETE("/projects/:id", deleteProjectController.Execute)
}
