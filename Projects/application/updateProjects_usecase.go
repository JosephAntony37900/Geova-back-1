package application

import (
	"fmt"
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

func NewUpdateProjectUseCase(
	repo repository.ProjectRepository,
	imageUploadSvc *services.ImageUploadWorkerService,
) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{
		repo:           repo,
		imageUploadSvc: imageUploadSvc,
	}
}

func (uc *UpdateProjectUseCase) Execute(project entities.Project, imagePath string) error {
	if imagePath != "" {
		log.Printf("INFO: [Update] Procesando imagen para proyecto ID=%d", project.Id)
		
		result, err := uc.imageUploadSvc.SubmitUploadJobSync(
			project.Id,
			imagePath,
			30*time.Second, 
		)

		if err != nil {
			return fmt.Errorf("error al subir imagen: %w", err)
		}

		if result.Error != nil {
			return fmt.Errorf("error al subir imagen: %w", result.Error)
		}

		log.Printf("SUCCESS: [Update] Imagen subida en %v: %s", result.Duration, result.ImageURL)
		project.Img = result.ImageURL
	}

	log.Printf("INFO: [Update] Actualizando proyecto ID=%d en base de datos", project.Id)
	return uc.repo.Update(&project)
}