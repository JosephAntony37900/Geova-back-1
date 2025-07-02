package application
import (

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type CreateProjectUseCase struct {
	db repository.ProjectRepository
}

func NewCreateProjectUseCase (db repository.ProjectRepository ) *CreateProjectUseCase{
	return &CreateProjectUseCase{
		db: db,
	}
}

func (uc *CreateProjectUseCase) Execute(Project entities.Project) error{
	return uc.db.Save(Project)
}