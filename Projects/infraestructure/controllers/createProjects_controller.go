package controllers


import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
)

type CreateProjectController struct {
	useCase *application.CreateProjectUseCase
}

func NewCreateProjectController(useCase *application.CreateProjectUseCase) *CreateProjectController {
	return &CreateProjectController{useCase: useCase}
}

func (c *CreateProjectController) Execute(ctx *gin.Context) {
	var project entities.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(400, gin.H{"error": "Error al leer los datos"})
		return
	}

	err := c.useCase.Execute(project)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error al crear proyecto: " + err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Proyecto creado exitosamente"})
}
