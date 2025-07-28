// application/update_user_usecase.go
package application

import (
	"fmt"
	"strings"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/services"
)

type UpdateUserUseCase struct {
	repo   repository.UserRepository
	bcrypt services.IBcryptService
}

func NewUpdateUserUseCase(repo repository.UserRepository, bcrypt services.IBcryptService) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (uc *UpdateUserUseCase) Execute(userUpdate entities.User) (*entities.User, error) {
	// Verificar que el usuario existe
	existingUser, err := uc.repo.FindById(userUpdate.Id)
	if err != nil {
		return nil, fmt.Errorf("usuario con ID %d no encontrado: %w", userUpdate.Id, err)
	}

	// Validaciones de negocio
	if err := uc.validateUpdateData(userUpdate); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Verificar que el email no esté siendo usado por otro usuario
	if userUpdate.Email != existingUser.Email {
		userWithEmail, _ := uc.repo.FindByEmail(userUpdate.Email)
		if userWithEmail != nil && userWithEmail.Id != userUpdate.Id {
			return nil, fmt.Errorf("el email %s ya está siendo usado por otro usuario", userUpdate.Email)
		}
	}

	// Preparar datos actualizados
	updatedUser := *existingUser
	updatedUser.Username = strings.TrimSpace(userUpdate.Username)
	updatedUser.Nombre = strings.TrimSpace(userUpdate.Nombre)
	updatedUser.Apellidos = strings.TrimSpace(userUpdate.Apellidos)
	updatedUser.Email = strings.ToLower(strings.TrimSpace(userUpdate.Email))

	// Solo actualizar contraseña si se proporciona una nueva
	if userUpdate.Password != "" {
		hashedPassword, err := uc.bcrypt.HashPassword(userUpdate.Password)
		if err != nil {
			return nil, fmt.Errorf("error al procesar la nueva contraseña: %w", err)
		}
		updatedUser.Password = hashedPassword
	}

	// Actualizar usuario (el repository maneja la sincronización automáticamente)
	if err := uc.repo.Update(updatedUser); err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	// Obtener usuario actualizado y limpiar contraseña para respuesta
	finalUser, err := uc.repo.FindById(updatedUser.Id)
	if err != nil {
		// Si no podemos recuperarlo, devolver lo que teníamos
		finalUser = &updatedUser
	}
	
	finalUser.Password = ""
	return finalUser, nil
}

func (uc *UpdateUserUseCase) validateUpdateData(user entities.User) error {
	if user.Id <= 0 {
		return fmt.Errorf("ID de usuario inválido")
	}

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

	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	// Solo validar contraseña si se está intentando cambiar
	if user.Password != "" && len(user.Password) < 6 {
		return fmt.Errorf("la nueva contraseña debe tener al menos 6 caracteres")
	}

	return nil
}

func (uc *UpdateUserUseCase) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}