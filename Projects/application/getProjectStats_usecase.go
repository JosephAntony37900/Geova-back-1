package application

import (
    "fmt"
    "log"

    "github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
    "github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type GetProjectStatsUseCase struct {
    db repository.ProjectRepository
}

func NewGetProjectStatsUseCase(db repository.ProjectRepository) *GetProjectStatsUseCase {
    return &GetProjectStatsUseCase{
        db: db,
    }
}

func (uc *GetProjectStatsUseCase) Execute(userId int, days int) (*entities.ProjectStats, error) {
    if userId <= 0 {
        return nil, fmt.Errorf("userId debe ser mayor a 0")
    }

    if days <= 0 {
        days = 7 
    }

    log.Printf("INFO: Obteniendo estadísticas de proyectos - UserId: %d, Días: %d", userId, days)

    dailyCounts, err := uc.db.GetProjectsStats(userId, days)
    if err != nil {
        log.Printf("ERROR: Error al obtener estadísticas: %v", err)
        return nil, err
    }

    totalCount := 0
    for _, dc := range dailyCounts {
        totalCount += dc.Count
    }

    stats := &entities.ProjectStats{
        UserId:     userId,
        TotalCount: totalCount,
        Daily:      dailyCounts,
    }

    log.Printf("SUCCESS: Estadísticas obtenidas - Total: %d proyectos en %d días", totalCount, len(dailyCounts))

    return stats, nil
}