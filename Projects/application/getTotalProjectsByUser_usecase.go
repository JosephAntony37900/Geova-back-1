package application

import (
	"fmt"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetTotalProjectsByUserUseCase struct {
	projectRepo repository.ProjectRepository
}	

func NewGetTotalProjectsByUserUseCase(repo repository.ProjectRepository) *GetTotalProjectsByUserUseCase {
	return &GetTotalProjectsByUserUseCase{projectRepo: repo}
}

func (uc *GetTotalProjectsByUserUseCase) Execute(userId string) (int, error) {
    if userId == "" {
        return 0, fmt.Errorf("el ID de usuario es requerido")
    }
    
    count, err := uc.projectRepo.GetTotalProjectsByUser(userId)
    if err != nil {
        return 0, fmt.Errorf("error al obtener total de proyectos: %w", err)
    }
    
    return count, nil
}