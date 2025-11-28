package controllers

import (
	"net/http"
	"regexp"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Users/application"
)

type LoginUserController struct {
	useCase *application.LoginUseCase
}

func NewLoginUserController(useCase *application.LoginUseCase) *LoginUserController {
	return &LoginUserController{useCase: useCase}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *LoginUserController) Execute(ctx *gin.Context) {
	var req loginRequest
	
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

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	output, err := c.useCase.Execute(application.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"token":   output.Token,
		"user": gin.H{
			"id":        output.User.Id,
			"nombre":    output.User.Nombre,
			"apellidos": output.User.Apellidos,
			"email":     output.User.Email,
		},
	})
}

func (c *LoginUserController) validateRequest(req *loginRequest) error {
	if req.Email == "" {
		return fmt.Errorf("el campo email es requerido")
	}
	if req.Password == "" {
		return fmt.Errorf("el campo password es requerido")
	}

	if !isValidEmailFormat(req.Email) {
		return fmt.Errorf("formato de email inválido")
	}

	if len(req.Password) < 8 {
		return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
	}

	if len(req.Password) > 128 {
		return fmt.Errorf("la contraseña es demasiado larga")
	}

	return nil
}

func isValidEmailFormat(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

