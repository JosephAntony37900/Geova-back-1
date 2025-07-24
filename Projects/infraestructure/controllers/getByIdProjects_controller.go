package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)

type GetProjectByIdController struct{
	useCase *application.GetProjectById
}

func NewGetProjectByIdUseController(usecase *application.GetProjectById) *GetProjectByIdController{
	return &GetProjectByIdController{useCase: usecase}
}

func (c *GetProjectByIdController) Execute(ctx *gin.Context){
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H {"error": "ID invalido"})
		return
	}

	project, err := c.useCase.Execute(id)
	if  err != nil {
		ctx.JSON(http.StatusNotFound, gin.H {"error": "Proyecto inexistente"})
		return
	}
	ctx.JSON(http.StatusOK, project)
}	