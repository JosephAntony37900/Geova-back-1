package repository

import (

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
)

type ProjectRepository interface {
	Save(project *entities.Project) error
	FindById(id int) (*entities.Project, error)
	FindAll() ([]entities.Project, error)
	Update(project *entities.Project) error
	Delete (id int) error
	FindByName(nombre string) ([]entities.Project, error)
	FindByCategory(categoria string) ([]entities.Project, error)
	FindByDate(fecha string) ([]entities.Project, error)
	FindByUserId(userId int) ([]entities.Project, error)
}
