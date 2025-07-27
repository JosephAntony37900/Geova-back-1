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
	db       repository.ProjectRepository
	cloudSrv services.ICloudinaryService
}

type ProjectCreationResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	IsOffline  bool   `json:"is_offline"`
	HasImage   bool   `json:"has_image"`
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

// Verificar si el error es relacionado con conectividad
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
			
			log.Println("INFO: Conectividad disponible, intentando subir imagen a Cloudinary...")
			
			url, err := uc.cloudSrv.UploadImage(imagePath)
			if err != nil {
				
				if uc.isConnectivityError(err) {
					log.Printf("WARNING: Error de conectividad detectado al subir imagen: %v", err)
					result.IsOffline = true
					
					
					project.Img = ""
					result.Message = "Proyecto creado sin imagen debido a problemas de conectividad. La imagen se subirá cuando haya conexión a internet."
				} else {
					
					log.Printf("ERROR: Error al subir imagen (no conectividad): %v", err)
					return result, err
				}
			} else {
				
				project.Img = url
				result.Message = "Proyecto creado exitosamente con imagen"
				log.Printf("SUCCESS: Imagen subida exitosamente: %s", url)
			}
		} else {
			
			log.Println("WARNING: Sin conectividad a internet, creando proyecto sin imagen")
			result.IsOffline = true
			project.Img = ""
			result.Message = "Proyecto creado sin imagen debido a falta de conexión a internet."
		}
	} else {
		
		project.Img = ""
		result.Message = "Proyecto creado exitosamente sin imagen"
		log.Println("INFO: Proyecto creado sin imagen (no se proporcionó archivo)")
	}

	
	if err := uc.db.Save(project); err != nil {
		log.Printf("ERROR: Error al guardar proyecto en BD: %v", err)
		return result, err
	}

	result.Success = true
	log.Printf("SUCCESS: Proyecto creado - ID: %d, Offline: %t, HasImage: %t", 
		project.Id, result.IsOffline, result.HasImage)
	
	return result, nil
}