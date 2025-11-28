// geova-back-1/Users/application/updateUsers_useCase.go
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

type UpdateUserInput struct {
	Id        int
	Username  string
	Nombre    string
	Apellidos string
	Email     string
	Password  string // Opcional - solo si se quiere cambiar
}

type UpdateUserOutput struct {
	User *entities.User
}

func (uc *UpdateUserUseCase) Execute(input UpdateUserInput) (*UpdateUserOutput, error) {
	// Validación de negocio: campos requeridos
	if err := uc.validateBusinessRules(input); err != nil {
		return nil, err
	}

	// Validación de negocio: el usuario existe
	existingUser, err := uc.repo.FindById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Validación de negocio: email único (si cambió)
	if input.Email != existingUser.Email {
		if err := uc.validateEmailUniqueness(input.Email, input.Id); err != nil {
			return nil, err
		}
	}

	// Construir el usuario actualizado
	updatedUser := uc.buildUpdatedUser(existingUser, input)

	// Validación de negocio: si se cambia contraseña, hashearla
	if input.Password != "" {
		hashedPassword, err := uc.bcrypt.HashPassword(input.Password)
		if err != nil {
			return nil, fmt.Errorf("error al procesar la contraseña")
		}
		updatedUser.Password = hashedPassword
	}

	// Persistir cambios
	if err := uc.repo.Update(*updatedUser); err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	// Obtener usuario actualizado
	finalUser, err := uc.repo.FindById(updatedUser.Id)
	if err != nil {
		finalUser = updatedUser
	}

	// Limpiar contraseña antes de retornar
	finalUser.Password = ""

	return &UpdateUserOutput{
		User: finalUser,
	}, nil
}

// validateBusinessRules valida las reglas de negocio básicas
func (uc *UpdateUserUseCase) validateBusinessRules(input UpdateUserInput) error {
	if input.Id <= 0 {
		return fmt.Errorf("ID de usuario inválido")
	}

	if strings.TrimSpace(input.Username) == "" {
		return fmt.Errorf("el nombre de usuario es requerido")
	}

	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if strings.TrimSpace(input.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	// Regla de negocio: si se proporciona contraseña, debe cumplir longitud mínima
	if input.Password != "" && len(input.Password) < 8 {
		return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
	}

	return nil
}

// validateEmailUniqueness valida que el email no esté en uso por otro usuario
func (uc *UpdateUserUseCase) validateEmailUniqueness(email string, userId int) error {
	userWithEmail, err := uc.repo.FindByEmail(email)
	if err != nil {
		// Si no existe usuario con ese email, está disponible
		return nil
	}

	// Si existe pero es el mismo usuario, está bien
	if userWithEmail.Id == userId {
		return nil
	}

	// Email en uso por otro usuario
	return fmt.Errorf("el email ya está siendo usado por otro usuario")
}

// buildUpdatedUser construye la entidad con los datos actualizados
func (uc *UpdateUserUseCase) buildUpdatedUser(existing *entities.User, input UpdateUserInput) *entities.User {
	updated := *existing
	updated.Username = input.Username
	updated.Nombre = input.Nombre
	updated.Apellidos = input.Apellidos
	updated.Email = input.Email
	
	return &updated
}