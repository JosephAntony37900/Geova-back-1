package repository

import (

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
)

type ProjectRepository interface {
	Save(proyect entities.Project) error
	FindById(id int) (*entities.Project, error)
	FindAll() ([]entities.Project, error)
	Update(proyect entities.Project) error
	Delete (id int) error
	FindByName(nombre string) ([]entities.Project, error)
	FindByCategory(categoria string) ([]entities.Project, error)
	FindByDate(fecha string) ([]entities.Project, error)
	FindByUserId(userId int) ([]entities.Project, error)
	GetProjectsStats(userId int, days int) ([]entities.DailyProjectCount, error)
	GetTotalProjectsByUser(userId string) (int, error)
}
