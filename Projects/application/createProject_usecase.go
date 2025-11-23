package application

import (
	"log"
	"net"
	"strings"
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
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	IsOffline bool   `json:"is_offline"`
	HasImage  bool   `json:"has_image"`
}

func NewCreateProjectUseCase(db repository.ProjectRepository, imageUploadSvc *services.ImageUploadWorkerService) *CreateProjectUseCase {
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

func (uc *CreateProjectUseCase) isConnectivityError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())
	connectivityKeywords := []string{
		"no such host",
		"dial tcp",
		"connection refused",
		"timeout",
		"network is unreachable",
		"temporary failure in name resolution",
	}

	for _, keyword := range connectivityKeywords {
		if strings.Contains(errorStr, keyword) {
			return true
		}
	}
	return false
}

func (uc *CreateProjectUseCase) Execute(project entities.Project, imagePath string) (*ProjectCreationResult, error) {
	result := &ProjectCreationResult{
		Success:   false,
		IsOffline: false,
		HasImage:  imagePath != "",
	}

	hasInternet := uc.hasInternetConnection()

	if imagePath != "" {
		if hasInternet {
			// Conectividad disponible, encolar trabajo de subida de imagen
			log.Println("INFO: Conectividad disponible, encolando trabajo de subida de imagen...")

			// Usar SubmitUploadJob (asíncrono) - el proyecto se guarda sin esperar
			err := uc.imageUploadSvc.SubmitUploadJob(project.Id, imagePath)
			if err != nil {
				// Si falla el encolamiento (servicio cerrado, cola llena), tratar como offline
				log.Printf("WARNING: No se pudo encolar trabajo de subida: %v", err)
				result.IsOffline = true
				project.Img = ""
				result.Message = "Proyecto creado sin imagen - servicio de subida no disponible. Se intentará más tarde."
			} else {
				// Trabajo encolado exitosamente, pero aún no tenemos URL
				// Guardar proyecto sin imagen por ahora
				project.Img = ""
				result.Message = "Proyecto creado exitosamente - imagen en proceso de subida"
				log.Printf("INFO: Trabajo de subida encolado para proyecto %d", project.Id)
			}
		} else {
			// Sin conectividad
			log.Println("WARNING: Sin conectividad a internet, creando proyecto sin imagen")
			result.IsOffline = true
			project.Img = ""
			result.Message = "Proyecto creado sin imagen debido a falta de conexión a internet."
		}
	} else {
		// Sin imagen
		project.Img = ""
		result.Message = "Proyecto creado exitosamente sin imagen"
		log.Println("INFO: Proyecto creado sin imagen (no se proporcionó archivo)")
	}

	// Guardar proyecto en BD
	if err := uc.db.Save(project); err != nil {
		log.Printf("ERROR: Error al guardar proyecto en BD: %v", err)
		return result, err
	}

	result.Success = true
	log.Printf("SUCCESS: Proyecto creado - ID: %d, Offline: %t, HasImage: %t",
		project.Id, result.IsOffline, result.HasImage)

	return result, nil
}
