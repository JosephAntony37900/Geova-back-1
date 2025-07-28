// application/create_user_usecase.go
package application

import (
	"fmt"
	"strings"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/services"
)

type CreateUserUseCase struct {
	repo   repository.UserRepository
	bcrypt services.IBcryptService
}

func NewCreateUserUseCase(repo repository.UserRepository, bcrypt services.IBcryptService) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (uc *CreateUserUseCase) Execute(user entities.User) (*entities.User, error) {
	// Validaciones de negocio
	if err := uc.validateUser(user); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Verificar que el email no exista (el repository ya maneja esto, pero es buena práctica validar aquí)
	existingUser, _ := uc.repo.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("el email %s ya está registrado", user.Email)
	}

	// Hashear la contraseña
	hashedPassword, err := uc.bcrypt.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("error al procesar la contraseña: %w", err)
	}
	user.Password = hashedPassword

	// Normalizar datos
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)
	user.Nombre = strings.TrimSpace(user.Nombre)
	user.Apellidos = strings.TrimSpace(user.Apellidos)

	// Guardar usuario (el repository maneja la sincronización automáticamente)
	if err := uc.repo.Save(user); err != nil {
		return nil, fmt.Errorf("error al guardar usuario: %w", err)
	}

	// Buscar el usuario creado para obtener el ID asignado
	createdUser, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		// Aunque se guardó, no pudimos recuperarlo - esto es raro pero manejable
		return &user, nil
	}

	// No devolver la contraseña hasheada en la respuesta
	createdUser.Password = ""
	
	return createdUser, nil
}

func (uc *CreateUserUseCase) validateUser(user entities.User) error {
	if strings.TrimSpace(user.Username) == "" {
		return fmt.Errorf("el nombre de usuario es requerido")
	}

	if len(user.Username) < 3 {
		return fmt.Errorf("el nombre de usuario debe tener al menos 3 caracteres")
	}

	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if !uc.isValidEmail(user.Email) {
		return fmt.Errorf("el formato del email no es válido")
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

func (uc *CreateUserUseCase) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}