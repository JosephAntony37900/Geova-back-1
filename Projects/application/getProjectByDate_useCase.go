package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetProjectsByDateUseCase struct {
	projectRepo repository.ProjectRepository
}

func NewGetProjectsByDateUseCase(repo repository.ProjectRepository) *GetProjectsByDateUseCase {
	return &GetProjectsByDateUseCase{projectRepo: repo}
}

func (uc *GetProjectsByDateUseCase) Execute(fecha string) ([]entities.Project, error) {
	return uc.projectRepo.FindByDate(fecha)
}
