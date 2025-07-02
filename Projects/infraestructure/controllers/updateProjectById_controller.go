package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)


type UpdateProjectController struct {
	useCase *application.UpdateProjectUseCase
}

func NewUpdateProjectController(useCase *application.UpdateProjectUseCase) *UpdateProjectController {
	return &UpdateProjectController{useCase: useCase}
}

func (c *UpdateProjectController) Execute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var project entities.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	project.Id = id
	if err := c.useCase.Execute(project); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la informacion del proyecto: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "La información del proyecto se actualizo correctamente"})
}
