package application

import (
	"fmt"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
)

type UpdateUserCase struct {
	db repository.UserRepository
}

func NewUpdateUserCase(db repository.UserRepository) *UpdateUserCase {
	return &UpdateUserCase{db: db}
}

func (uu *UpdateUserCase) Execute(userId int, userData map[string]interface{}) error {
	user, err := uu.db.FindById(userId)
	if err != nil {
		return err
	}

	if nombre, ok := userData["nombre"].(string); ok {
		user.Nombre = nombre
	}
	if apellidos, ok := userData["apellidos"].(string); ok {
		user.Apellidos = apellidos
	}
	if email, ok := userData["email"].(string); ok {
		user.Email = email
	}
	if password, ok := userData["password"].(string); ok {
		user.Password = password
	}
    // guardo los cambios en el repositorio
	if err := uu.db.Update(*user); err != nil {
		return fmt.Errorf("error actualizando el usuario: %w", err)
	}

	return nil
}