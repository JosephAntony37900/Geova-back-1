package controllers

import (
	"net/http"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)

type GetProjectByDateController struct {
	useCase *application.GetProjectsByDateUseCase
}

func NewGetProjectByDateController(usecase *application.GetProjectsByDateUseCase) *GetProjectByDateController {
	return &GetProjectByDateController{useCase: usecase}
}

func (c *GetProjectByDateController) Execute(ctx *gin.Context) {
	nombre := ctx.Param("fecha")

	projects, err := c.useCase.Execute(nombre)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}
