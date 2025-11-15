// controllers/update_user_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Users/application"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
)

type UpdateUserController struct {
	useCase *application.UpdateUserUseCase
}

func NewUpdateUserController(useCase *application.UpdateUserUseCase) *UpdateUserController {
	return &UpdateUserController{useCase: useCase}
}

func (c *UpdateUserController) Execute(ctx *gin.Context) {
	// Validar y obtener ID del parámetro
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
			"details": "El ID debe ser un número entero positivo",
		})
		return
	}

	// Bindear datos JSON
	var userUpdate entities.User
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al leer los datos del usuario",
			"details": err.Error(),
		})
		return
	}

	// Asignar ID del parámetro URL
	userUpdate.Id = id

	// Validaciones básicas
	if err := c.validateUpdateInput(userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de actualización inválidos",
			"details": err.Error(),
		})
		return
	}

	// Ejecutar caso de uso de actualización
	updatedUser, err := c.useCase.Execute(userUpdate)
	if err != nil {
		// Clasificar errores
		if strings.Contains(err.Error(), "no existe") || 
		   strings.Contains(err.Error(), "no encontrado") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
				"details": err.Error(),
			})
			return
		}
		
		if strings.Contains(err.Error(), "ya está siendo usado") || 
		   strings.Contains(err.Error(), "email duplicado") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Email ya está en uso",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno al actualizar usuario",
			"details": err.Error(),
		})
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
		"user": gin.H{
			"id": updatedUser.Id,
			"username": updatedUser.Username,
			"email": updatedUser.Email,
			"nombre": updatedUser.Nombre,
			"apellidos": updatedUser.Apellidos,
		},
		"sync_status": "local_updated",
	})
}

func (c *UpdateUserController) validateUpdateInput(user entities.User) error {
	if strings.TrimSpace(user.Username) == "" {
		return fmt.Errorf("el nombre de usuario es requerido")
	}
	
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}
	
	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}
	
	// Solo validar password si se está intentando cambiar
	if user.Password != "" && len(user.Password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}
	
	return nil
}