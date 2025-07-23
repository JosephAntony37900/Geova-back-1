package controllers

import (
	"fmt"
	"net/http"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"

)

type GetProjectByDateController struct {
	useCase *application.GetProjectsByDateUseCase
}

func NewGetProjectByDateController(usecase *application.GetProjectsByDateUseCase) *GetProjectByDateController {
	return &GetProjectByDateController{useCase: usecase}
}

func (c *GetProjectByDateController) Execute(ctx *gin.Context) {
	fecha := ctx.Param("fecha") 
	
	
	fmt.Printf("DEBUG GetProjectByDate - Fecha recibida: '%s'\n", fecha)

	projects, err := c.useCase.Execute(fecha)
	if err != nil {
		fmt.Printf("DEBUG GetProjectByDate - Error: %v\n", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("DEBUG GetProjectByDate - Proyectos encontrados: %d\n", len(projects))
	
	
	if len(projects) == 0 {
		ctx.JSON(http.StatusOK, []entities.Project{})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}