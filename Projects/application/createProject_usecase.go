package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type CreateProjectUseCase struct {
	db       repository.ProjectRepository
	cloudSrv services.ICloudinaryService
}

func NewCreateProjectUseCase(db repository.ProjectRepository, cloudSrv services.ICloudinaryService) *CreateProjectUseCase {
	return &CreateProjectUseCase{
		db:       db,
		cloudSrv: cloudSrv,
	}
}

func (uc *CreateProjectUseCase) Execute(project entities.Project, imagePath string) error {
	url, err := uc.cloudSrv.UploadImage(imagePath)
	if err != nil {
		return err
	}
	project.Img = url
	return uc.db.Save(project)
}
