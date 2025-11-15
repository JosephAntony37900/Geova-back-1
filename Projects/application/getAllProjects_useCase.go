package application


import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetAllProjectsUseCase struct {
	db repository.ProjectRepository
}

func NewGeProjectsUseCase(db repository.ProjectRepository) *GetAllProjectsUseCase {
	return &GetAllProjectsUseCase{db: db}
}

func (gp *GetAllProjectsUseCase) Execute() ([]entities.Project, error) {
	project, err := gp.db.FindAll()
	if err != nil {
		return nil, err
	}
	return project, nil
}