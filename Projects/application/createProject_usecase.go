package application

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

type CreateProjectUseCase struct {
	db       repository.ProjectRepository
	cloudSrv services.ICloudinaryService
	mu       sync.Mutex 
}

type ProjectCreationResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	IsOffline  bool   `json:"is_offline"`
	HasImage   bool   `json:"has_image"`
	ProjectID  int    `json:"project_id,omitempty"` 
}

func NewCreateProjectUseCase(db repository.ProjectRepository, cloudSrv services.ICloudinaryService) *CreateProjectUseCase {
	return &CreateProjectUseCase{
		db:       db,
		cloudSrv: cloudSrv,
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

// uploadImageAsync sube la imagen en una goroutine y actualiza el proyecto tambien se hace una img temporal
func (uc *CreateProjectUseCase) uploadImageAsync(projectID int, imagePath string) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		
		log.Printf("INFO: [Async] Iniciando subida de imagen para proyecto ID=%d", projectID)
		
		if projectID <= 0 {
			log.Printf("ERROR: [Async] ID de proyecto inválido: %d", projectID)
			return
		}
		
		url, err := uc.cloudSrv.UploadImage(imagePath)
		if err != nil {
			log.Printf("ERROR: [Async] Error al subir imagen para proyecto ID=%d: %v", projectID, err)
			return
		}

		log.Printf("SUCCESS: [Async] Imagen subida exitosamente: %s", url)

		
		project, err := uc.db.FindById(projectID)
		if err != nil {
			log.Printf("ERROR: [Async] No se pudo encontrar proyecto ID=%d para actualizar imagen: %v", projectID, err)
			return
		}

		project.Img = url
		
		uc.mu.Lock()
		err = uc.db.Update(project) 
		uc.mu.Unlock()

		if err != nil {
			log.Printf("ERROR: [Async] Error al actualizar proyecto ID=%d con URL de imagen: %v", projectID, err)
			return
		}

		log.Printf("SUCCESS: [Async] Proyecto ID=%d actualizado con imagen: %s", projectID, url)
	}()
}

// primero se guarda el proyecto en la BD sin imagen, luego si hay imagen se sube en segundo plano
func (uc *CreateProjectUseCase) Execute(project entities.Project, imagePath string) (*ProjectCreationResult, error) {
	result := &ProjectCreationResult{
		Success:   false,
		IsOffline: false,
		HasImage:  imagePath != "",
	}

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
			log.Println("INFO: Lanzando subida asíncrona de imagen...")
			uc.uploadImageAsync(project.Id, imagePath)
			
			result.Message = "Proyecto creado exitosamente. La imagen se está subiendo en segundo plano."
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