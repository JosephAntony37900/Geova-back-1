package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntony37900/Geova-back-1/Users/application"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
)

type CreateUserController struct {
	useCase *application.CreateUserUseCase
}

func NewCreateUserController(useCase *application.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{useCase: useCase}
}

func (c *CreateUserController) Execute(ctx *gin.Context) {
	var user entities.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": "Error al leer los datos"})
		return
	}

	err := c.useCase.Execute(user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error al crear usuario: " + err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Usuario creado exitosamente"})
}
