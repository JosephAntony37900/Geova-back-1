package controllers

import (
	"net/http"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)

type GetProjectByNameController struct {
	useCase *application.GetProjectsByNameUseCase
}

func NewGetProjectByNameController(usecase *application.GetProjectsByNameUseCase) *GetProjectByNameController {
	return &GetProjectByNameController{useCase: usecase}
}

func (c *GetProjectByNameController) Execute(ctx *gin.Context) {
	nombre := ctx.Param("nombre")

	projects, err := c.useCase.Execute(nombre)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}
