package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
)

type SyncProjectsController struct {
	useCase *application.SyncProjectsUseCase
}

func NewSyncProjectsController(useCase *application.SyncProjectsUseCase) *SyncProjectsController {
	return &SyncProjectsController{useCase: useCase}
}

func (c *SyncProjectsController) Execute(ctx *gin.Context) {
	var projects []entities.Project
	if err := ctx.ShouldBindJSON(&projects); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if err := c.useCase.Execute(projects); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al sincronizar proyectos"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Proyectos sincronizados correctamente"})
}
