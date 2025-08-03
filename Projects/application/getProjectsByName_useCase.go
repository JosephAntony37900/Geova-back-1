package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetProjectsByNameUseCase struct {
	projectRepo repository.ProjectRepository
}

func NewGetProjectsByNameUseCase(repo repository.ProjectRepository) *GetProjectsByNameUseCase {
	return &GetProjectsByNameUseCase{projectRepo: repo}
}

func (uc *GetProjectsByNameUseCase) Execute(nombre string) ([]entities.Project, error) {
	return uc.projectRepo.FindByName(nombre)
}
