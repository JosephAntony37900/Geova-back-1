package repository

import "github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"

type ProjectRepository interface {
	Save(proyect entities.Project) error
	FindById(id int) (*entities.Project, error)
	FindAll() ([]entities.Project, error)
	Update(proyect entities.Project) error
	Delete (id int) error
}
