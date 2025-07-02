package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)


type UpdateProjectUseCase struct {
	repo repository.ProjectRepository
}

func NewUpdateProjectUseCase(repo repository.ProjectRepository) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{repo: repo}
}

func (up *UpdateProjectUseCase) Execute(project entities.Project) error{
	return up.repo.Update(project)
}