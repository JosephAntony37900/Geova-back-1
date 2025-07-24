package controllers

import (
	"net/http"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)

type GetProjectByCategoryController struct {
	useCase *application.GetProjectsByCategoryUseCase
}

func NewGetProjectByCategoryController(usecase *application.GetProjectsByCategoryUseCase) *GetProjectByCategoryController {
	return &GetProjectByCategoryController{useCase: usecase}
}

func (c *GetProjectByCategoryController) Execute(ctx *gin.Context) {
	nombre := ctx.Param("categoria")

	projects, err := c.useCase.Execute(nombre)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}
