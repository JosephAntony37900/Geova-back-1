package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/gin-gonic/gin"
)

type GetProjectsByUserIdController struct {
	useCase *application.GetProjectsByUserId
}

func NewGetProjectsByUserIdController(usecase *application.GetProjectsByUserId) *GetProjectsByUserIdController {
	return &GetProjectsByUserIdController{useCase: usecase}
}

func (c *GetProjectsByUserIdController) Execute(ctx *gin.Context) {
	userIdParam := ctx.Param("userId")
	userId, err := strconv.Atoi(userIdParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	projects, err := c.useCase.Execute(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron proyectos para este usuario"})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}