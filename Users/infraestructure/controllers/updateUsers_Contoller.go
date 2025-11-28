// geova-back-1/Users/controllers/updateUsers_Controller.go
package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Users/application"
)

type UpdateUserController struct {
	useCase *application.UpdateUserUseCase
}

func NewUpdateUserController(useCase *application.UpdateUserUseCase) *UpdateUserController {
	return &UpdateUserController{useCase: useCase}
}

type updateUserRequest struct {
	Username  string `json:"username"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"` // Opcional
}

func (c *UpdateUserController) Execute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido en la URL",
		})
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de datos inválido",
		})
		return
	}

	if err := c.validateRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	input := application.UpdateUserInput{
		Id:        id,
		Username:  strings.TrimSpace(req.Username),
		Nombre:    strings.TrimSpace(req.Nombre),
		Apellidos: strings.TrimSpace(req.Apellidos),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  strings.TrimSpace(req.Password),
	}

	output, err := c.useCase.Execute(input)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
		"user": gin.H{
			"id":        output.User.Id,
			"username":  output.User.Username,
			"email":     output.User.Email,
			"nombre":    output.User.Nombre,
			"apellidos": output.User.Apellidos,
		},
		"sync_status": "local_updated",
	})
}

func (c *UpdateUserController) validateRequest(req *updateUserRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("el campo username es requerido")
	}
	if strings.TrimSpace(req.Nombre) == "" {
		return fmt.Errorf("el campo nombre es requerido")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("el campo email es requerido")
	}

	// Validación técnica: longitud mínima de username
	if len(strings.TrimSpace(req.Username)) < 3 {
		return fmt.Errorf("el username debe tener al menos 3 caracteres")
	}

	// Validación técnica: formato de email
	if err := c.validateEmail(req.Email); err != nil {
		return err
	}

	// Validación técnica: si se proporciona password, validar complejidad
	if req.Password != "" {
		if err := c.validatePassword(req.Password); err != nil {
			return err
		}
		if len(req.Password) > 128 {
			return fmt.Errorf("la contraseña no puede exceder 128 caracteres")
		}
	}

	return nil
}

// validateEmail valida el formato del correo electrónico
func (c *UpdateUserController) validateEmail(email string) error {
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

// validatePassword valida que la contraseña cumpla con los requisitos de seguridad
func (c *UpdateUserController) validatePassword(password string) error {
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

func (c *UpdateUserController) handleError(ctx *gin.Context, err error) {
	errorMsg := err.Error()

	if strings.Contains(errorMsg, "no encontrado") ||
		strings.Contains(errorMsg, "no existe") {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Usuario no encontrado",
			"details": errorMsg,
		})
		return
	}

	if strings.Contains(errorMsg, "ya está siendo usado") ||
		strings.Contains(errorMsg, "email duplicado") {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "Email ya está en uso",
			"details": errorMsg,
		})
		return
	}

	if strings.Contains(errorMsg, "requerido") ||
		strings.Contains(errorMsg, "debe tener al menos") ||
		strings.Contains(errorMsg, "inválido") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de actualización inválidos",
			"details": errorMsg,
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error":   "Error interno al actualizar usuario",
		"details": errorMsg,
	})
}