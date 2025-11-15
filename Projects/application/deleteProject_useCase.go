package application

import (
	"fmt"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type DeleleProjectUseCase struct {
	db repository.ProjectRepository
}

func NewDeleteProjectUseCase (db repository.ProjectRepository) *DeleleProjectUseCase{
	return &DeleleProjectUseCase{db: db}
}

func (dp *DeleleProjectUseCase) Execute(id int) error{
	_, err := dp.db.FindById(id)
	if err != nil {
		return fmt.Errorf("Proyecto con el ID %d no encontrado: %w", id, err )
	}
	if err := dp.db.Delete(id); err != nil{
		return fmt.Errorf("Error al eliminar el usuarios con ese ID %d: %w", id, err)
	}
	return nil
}