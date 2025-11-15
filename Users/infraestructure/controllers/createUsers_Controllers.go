// controllers/create_user_controller.go
package controllers

import (
	"net/http"
	"strings"
	"fmt"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al leer los datos del usuario",
			"details": err.Error(),
		})
		return
	}

	// Validaciones básicas
	if err := c.validateUserInput(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de usuario inválidos",
			"details": err.Error(),
		})
		return
	}

	// Ejecutar caso de uso
	createdUser, err := c.useCase.Execute(user)
	if err != nil {
		// Clasificar errores para respuestas más específicas
		if strings.Contains(err.Error(), "ya está registrado") || 
		   strings.Contains(err.Error(), "ya existe") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Email ya registrado",
				"details": err.Error(),
			})
			return
		}
		
		if strings.Contains(err.Error(), "validación") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Error de validación",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno al crear usuario",
			"details": err.Error(),
		})
		return
	}

	// Respuesta exitosa con información del usuario creado
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente",
		"user": gin.H{
			"id": createdUser.Id,
			"username": createdUser.Username,
			"email": createdUser.Email,
			"nombre": createdUser.Nombre,
			"apellidos": createdUser.Apellidos,
		},
		"sync_status": "local_saved", // Indica que se guardó localmente
	})
}

func (c *CreateUserController) validateUserInput(user entities.User) error {
	if strings.TrimSpace(user.Username) == "" {
		return fmt.Errorf("el nombre de usuario es requerido")
	}
	
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}
	
	if strings.TrimSpace(user.Password) == "" {
		return fmt.Errorf("la contraseña es requerida")
	}
	
	if len(user.Password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}
	
	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}
	
	return nil
}