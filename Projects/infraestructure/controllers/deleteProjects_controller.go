package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
)
type DeleteProjectController struct {
	useCase *application.DeleleProjectUseCase
}

func NewDeleteProjectController(useCase *application.DeleleProjectUseCase) *DeleteProjectController{
	return &DeleteProjectController{useCase: useCase}
}

func (c *DeleteProjectController) Execute(ctx *gin.Context){
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := c.useCase.Execute(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al eliminar proyecto"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Proyecto eliminado correctamente"})
}