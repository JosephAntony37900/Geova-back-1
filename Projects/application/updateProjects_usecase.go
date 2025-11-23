package application

import (
	"log"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type UpdateProjectUseCase struct {
	repo           repository.ProjectRepository
	imageUploadSvc *services.ImageUploadWorkerService
}

func NewUpdateProjectUseCase(repo repository.ProjectRepository, imageUploadSvc *services.ImageUploadWorkerService) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{
		repo:           repo,
		imageUploadSvc: imageUploadSvc,
	}
}

func (uc *UpdateProjectUseCase) Execute(project entities.Project, imagePath string) error {
	if imagePath != "" {
		log.Printf("INFO: Actualizando proyecto %d con nueva imagen: %s", project.Id, imagePath)

		// Usar SubmitUploadJobSync con timeout de 30 segundos
		result, err := uc.imageUploadSvc.SubmitUploadJobSync(project.Id, imagePath, 30*time.Second)
		if err != nil {
			log.Printf("ERROR: Error al subir imagen para proyecto %d: %v", project.Id, err)
			return err
		}

		project.Img = result.ImageURL
		log.Printf("SUCCESS: Imagen subida exitosamente para proyecto %d: %s (duración: %v)",
			project.Id, result.ImageURL, result.Duration)
	}

	return uc.repo.Update(project)
}
