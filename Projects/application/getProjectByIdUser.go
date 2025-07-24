package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetProjectsByUserId struct {
	db repository.ProjectRepository
}

func NewGetProjectsByUserIdUseCase(db repository.ProjectRepository) *GetProjectsByUserId {
	return &GetProjectsByUserId{db: db}
}

func (gpbui *GetProjectsByUserId) Execute(userId int) ([]entities.Project, error) {
	projects, err := gpbui.db.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return projects, nil
}