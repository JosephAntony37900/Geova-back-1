package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Users/application"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
)

type UpdateUserController struct {
	useCase *application.UpdateUserCase
}

func NewUpdateUserController(useCase *application.UpdateUserCase) *UpdateUserController {
	return &UpdateUserController{useCase: useCase}
}

func (c *UpdateUserController) Execute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var user entities.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	user.Id = id
	if err := c.useCase.Execute(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado correctamente"})
}
