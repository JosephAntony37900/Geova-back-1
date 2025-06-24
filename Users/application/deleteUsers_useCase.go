package application

import (
	"fmt"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
)

type DeleteUserUseCase struct {
	db repository.UserRepository
}

func NewDeleteUserUseCase(db repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{db: db}
}

func (du *DeleteUserUseCase) Execute(id int) error {
	_, err := du.db.FindById(id)
	if err != nil {
		return fmt.Errorf("usuario con id %d no encontrado: %w", id, err)
	}

	if err := du.db.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar el usuario con id %d: %w", id, err)
	}
	return nil
}