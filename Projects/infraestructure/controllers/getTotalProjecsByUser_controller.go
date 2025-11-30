package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/JosephAntony37900/Geova-back-1/Projects/application"
)

type GetTotalProjectsByUserController struct {
    useCase *application.GetTotalProjectsByUserUseCase
}

func NewGetTotalProjectsByUserController(useCase *application.GetTotalProjectsByUserUseCase) * GetTotalProjectsByUserController {
    return & GetTotalProjectsByUserController{
        useCase: useCase,
    }
}

// Execute maneja la petición para obtener el total de proyectos de un usuario
func (c * GetTotalProjectsByUserController) Execute(ctx *gin.Context) {
    userId := ctx.Param("userId")
    
    if userId == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "El ID de usuario es requerido",
        })
        return
    }
    
    totalProjects, err := c.useCase.Execute(userId)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error al obtener total de proyectos",
            "details": err.Error(),
        })
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "user_id": userId,
        "total_projects": totalProjects,
    })
}