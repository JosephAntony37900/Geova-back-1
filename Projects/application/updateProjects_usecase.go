package application

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type UpdateProjectUseCase struct {
	repo     repository.ProjectRepository
	cloudSrv services.ICloudinaryService
}

type UpdateResult struct {
	ImageURL string
	Error    error
}

func NewUpdateProjectUseCase(repo repository.ProjectRepository, cloudSrv services.ICloudinaryService) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{
		repo:     repo,
		cloudSrv: cloudSrv,
	}
}

func (uc *UpdateProjectUseCase) Execute(project entities.Project, imagePath string) error {
	if imagePath != "" {
		resultChan := make(chan UpdateResult, 1)
		
		// WaitGroup para asegurar que la goroutine termine
		var wg sync.WaitGroup
		wg.Add(1)

		// Goroutine para subir imagen
		go func() {
			defer wg.Done()
			
			log.Printf("INFO: [Update] Iniciando subida de imagen para proyecto ID=%d", project.Id)
			startTime := time.Now()
			
			url, err := uc.cloudSrv.UploadImage(imagePath)
			
			elapsed := time.Since(startTime)
			if err != nil {
				log.Printf("ERROR: [Update] Error al subir imagen (tiempo: %v): %v", elapsed, err)
				resultChan <- UpdateResult{Error: err}
				return
			}
			
			log.Printf("SUCCESS: [Update] Imagen subida en %v: %s", elapsed, url)
			resultChan <- UpdateResult{ImageURL: url, Error: nil}
		}()

		select {
		case result := <-resultChan:
			if result.Error != nil {
				return fmt.Errorf("error al subir imagen: %w", result.Error)
			}
			project.Img = result.ImageURL
		case <-time.After(30 * time.Second):
			return fmt.Errorf("timeout al subir imagen a Cloudinary")
		}

		wg.Wait()
		close(resultChan)
	}

	log.Printf("INFO: [Update] Actualizando proyecto ID=%d en base de datos", project.Id)
	return uc.repo.Update(&project)
}