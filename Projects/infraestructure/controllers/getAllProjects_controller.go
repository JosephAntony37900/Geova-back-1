package controllers

import (
	"net/http"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)

type GetAllProjectsController struct {
	useCase *application.GetAllProjectsUseCase

}

func NewGetAllProjectsController(useCase *application.GetAllProjectsUseCase) *GetAllProjectsController {
	return &GetAllProjectsController{useCase: useCase}
}

func (c *GetAllProjectsController) Execute(ctx *gin.Context) {
	proyects, err := c.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener al obtener la lista de proyectos: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, proyects)
}
