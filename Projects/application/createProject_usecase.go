package application

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type CreateProjectUseCase struct {
	db             repository.ProjectRepository
	imageUploadSvc *services.ImageUploadWorkerService 
}

type ProjectCreationResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	IsOffline  bool   `json:"is_offline"`
	HasImage   bool   `json:"has_image"`
	ProjectID  int    `json:"project_id,omitempty"`
}

func NewCreateProjectUseCase(
	db repository.ProjectRepository,
	imageUploadSvc *services.ImageUploadWorkerService,
) *CreateProjectUseCase {
	return &CreateProjectUseCase{
		db:             db,
		imageUploadSvc: imageUploadSvc,
	}
}

func (uc *CreateProjectUseCase) hasInternetConnection() bool {
	timeout := time.Duration(5 * time.Second)
	_, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	return err == nil
}



func (uc *CreateProjectUseCase) Execute(project entities.Project, imagePath string) (*ProjectCreationResult, error) {
	result := &ProjectCreationResult{
		Success:   false,
		IsOffline: false,
		HasImage:  imagePath != "",
	}

	// Guardar proyecto  en BD sin imagen inicialmente
	project.Img = ""
	
	if err := uc.db.Save(&project); err != nil {
		log.Printf("ERROR: Error al guardar proyecto en BD: %v", err)
		return result, err
	}

	if project.Id <= 0 {
		log.Printf("ERROR: El proyecto no recibió un ID válido de la BD. ID=%d", project.Id)
		return result, fmt.Errorf("error: la base de datos no retornó un ID válido para el proyecto")
	}

	result.Success = true
	result.ProjectID = project.Id
	log.Printf("SUCCESS: Proyecto creado - ID: %d", project.Id)

	if imagePath != "" {
		hasInternet := uc.hasInternetConnection()
		
		if hasInternet {
			err := uc.imageUploadSvc.SubmitUploadJob(project.Id, imagePath)
			if err != nil {
				log.Printf("WARNING: No se pudo encolar trabajo de imagen: %v", err)
				result.Message = "Proyecto creado, pero la imagen no pudo ser procesada"
			} else {
				result.Message = "Proyecto creado exitosamente. La imagen se está subiendo en segundo plano."
			}
			result.IsOffline = false
		} else {
			log.Println("WARNING: Sin conectividad a internet, imagen pendiente")
			result.IsOffline = true
			result.Message = "Proyecto creado sin imagen debido a falta de conexión a internet."
		}
	} else {
		result.Message = "Proyecto creado exitosamente sin imagen"
		log.Println("INFO: Proyecto creado sin imagen (no se proporcionó archivo)")
	}

	return result, nil
}