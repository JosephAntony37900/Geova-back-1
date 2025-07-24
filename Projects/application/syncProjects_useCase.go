package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type SyncProjectsUseCase struct {
	repo repository.ProjectRepository
}

func NewSyncProjectsUseCase(repo repository.ProjectRepository) *SyncProjectsUseCase {
	return &SyncProjectsUseCase{repo: repo}
}

func (uc *SyncProjectsUseCase) Execute(projects []entities.Project) error {
	return uc.repo.SaveManyProjects(projects)
}
