package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/application"
)

type SyncUsersController struct {
	useCase *application.SyncUsersUseCase
}

func NewSyncUsersController(useCase *application.SyncUsersUseCase) *SyncUsersController {
	return &SyncUsersController{useCase: useCase}
}

func (c *SyncUsersController) Execute(ctx *gin.Context) {
	var users []entities.User
	if err := ctx.ShouldBindJSON(&users); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if err := c.useCase.Execute(users); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al sincronizar usuarios"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Usuarios sincronizados correctamente"})
}
