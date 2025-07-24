package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type UpdateProjectUseCase struct {
	repo     repository.ProjectRepository
	cloudSrv services.ICloudinaryService
}

func NewUpdateProjectUseCase(repo repository.ProjectRepository, cloudSrv services.ICloudinaryService) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{
		repo:     repo,
		cloudSrv: cloudSrv,
	}
}

func (uc *UpdateProjectUseCase) Execute(project entities.Project, imagePath *string) error {
	current, err := uc.repo.FindById(project.Id)
	if err != nil {
		return err
	}
	if imagePath != nil && *imagePath != "" {
		url, err := uc.cloudSrv.UploadImage(*imagePath)
		if err != nil {
			return err
		}
		project.Img = url
	} else {
		project.Img = current.Img
	}

	return uc.repo.Update(project)
}
