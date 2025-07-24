package application

import (

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetProjectById struct {
	db repository.ProjectRepository
}

func NewGetProjectByIdUseCase (db repository.ProjectRepository) *GetProjectById{
	return &GetProjectById{db: db}
}

func (gpbi *GetProjectById) Execute (id int) (*entities.Project, error){
	Project, err := gpbi.db.FindById(id)
	if err != nil {
		return nil, err
	}
	return Project, nil
}