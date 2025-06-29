package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Users/application"
)
type DeleteUserController struct {
	useCase *application.DeleteUserUseCase
}

func NewDeleteUserController(useCase *application.DeleteUserUseCase) *DeleteUserController{
	return &DeleteUserController{useCase: useCase}
}

func (c *DeleteUserController) Execute(ctx *gin.Context){
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := c.useCase.Execute(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al eliminar usuario"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}