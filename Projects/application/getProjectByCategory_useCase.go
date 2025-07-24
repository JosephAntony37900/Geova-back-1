package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetProjectsByCategoryUseCase struct {
	projectRepo repository.ProjectRepository
}

func NewGetProjectsByCategoryUseCase(repo repository.ProjectRepository) *GetProjectsByCategoryUseCase {
	return &GetProjectsByCategoryUseCase{projectRepo: repo}
}

func (uc *GetProjectsByCategoryUseCase) Execute(categoria string) ([]entities.Project, error) {
	return uc.projectRepo.FindByCategory(categoria)
}
