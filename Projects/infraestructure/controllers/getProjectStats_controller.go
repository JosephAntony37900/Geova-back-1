package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/JosephAntony37900/Geova-back-1/Projects/application"
)

type GetProjectStatsController struct {
    useCase *application.GetProjectStatsUseCase
}

func NewGetProjectStatsController(useCase *application.GetProjectStatsUseCase) *GetProjectStatsController {
    return &GetProjectStatsController{useCase: useCase}
}

func (c *GetProjectStatsController) Execute(ctx *gin.Context) {
    // Obtener userId del query parameter
    userIdStr := ctx.Query("userId")
    if userIdStr == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error":   "El parámetro userId es obligatorio",
            "success": false,
        })
        return
    }

    userId, err := strconv.Atoi(userIdStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error":   "El userId debe ser un número válido",
            "success": false,
        })
        return
    }

    if userId <= 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error":   "El userId debe ser mayor a 0",
            "success": false,
        })
        return
    }

    // Obtener días (opcional, por defecto 7)
    days := 7
    daysStr := ctx.Query("days")
    if daysStr != "" {
        parsedDays, err := strconv.Atoi(daysStr)
        if err == nil && parsedDays > 0 {
            days = parsedDays
        }
    }

    // Ejecutar use case
    stats, err := c.useCase.Execute(userId, days)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Error al obtener estadísticas: " + err.Error(),
            "success": false,
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    stats,
    })
}