package application

import (
	"sync"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type UpdateProjectUseCase struct {
	repo      repository.ProjectRepository
	cloudSrv  services.ICloudinaryService
	workerSrv *services.ImageUploadWorkerService
	mu        sync.Mutex // Protege acceso a repo.Update
}

func NewUpdateProjectUseCase(repo repository.ProjectRepository, cloudSrv services.ICloudinaryService, workerSrv *services.ImageUploadWorkerService) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{
		repo:      repo,
		cloudSrv:  cloudSrv,
		workerSrv: workerSrv,
	}
}

func (uc *UpdateProjectUseCase) Execute(project entities.Project, imagePath string) error {
	if imagePath != "" {
		// Usar el worker service con timeout de 30 segundos
		url, err := uc.workerSrv.SubmitUploadJobSync(imagePath, 30*time.Second)
		if err != nil {
			return err
		}
		project.Img = url
	}

	// Proteger acceso concurrente al repositorio
	uc.mu.Lock()
	defer uc.mu.Unlock()

	return uc.repo.Update(project)
}
