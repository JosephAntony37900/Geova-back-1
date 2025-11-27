// controllers/create_user_controller.go
package controllers

import (
	"net/http"
	"strings"
	"fmt"
	"regexp"
	"unicode"
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

	// Validaciones mejoradas
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
		"sync_status": "local_saved",
	})
}

func (c *CreateUserController) validateUserInput(user entities.User) error {
	// Validar username
	if strings.TrimSpace(user.Username) == "" {
		return fmt.Errorf("el nombre de usuario es requerido")
	}
	
	// Validar email (formato y requerido)
	if err := c.validateEmail(user.Email); err != nil {
		return err
	}
	
	// Validar contraseña con los mismos requisitos del login
	if err := c.validatePassword(user.Password); err != nil {
		return err
	}
	
	// Validar nombre
	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}
	
	return nil
}

// validateEmail valida el formato del correo electrónico
func (c *CreateUserController) validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("el correo electrónico es requerido")
	}

	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)

	if err != nil {
		return fmt.Errorf("error validando el correo electrónico")
	}

	if !matched {
		return fmt.Errorf("el formato del correo electrónico no es válido")
	}

	return nil
}

// Funcion para validar la contraseña
func (c *CreateUserController) validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("la contraseña es requerida")
	}

	if len(password) < 8 {
		return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
	}

	var (
		hasUpper   bool
		hasNumber  bool
		hasSpecial bool
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/"

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
		for _, special := range specialChars {
			if char == special {
				hasSpecial = true
				break
			}
		}
	}

	if !hasUpper && !hasNumber && !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula, un número y un carácter especial")
	}

	if !hasUpper && !hasNumber {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula y un número")
	}

	if !hasUpper && !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula y un carácter especial")
	}

	if !hasNumber && !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos un número y un carácter especial")
	}

	if !hasUpper {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula")
	}

	if !hasNumber {
		return fmt.Errorf("la contraseña debe contener al menos un número")
	}

	if !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos un carácter especial (!@#$%^&*()_+-=[]{}|;:,.<>?/)")
	}

	return nil
}