package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
)

type UpdateProjectController struct {
	useCase *application.UpdateProjectUseCase
}

func NewUpdateProjectController(useCase *application.UpdateProjectUseCase) *UpdateProjectController {
	return &UpdateProjectController{useCase: useCase}
}

func (c *UpdateProjectController) Execute(ctx *gin.Context) {
	// Obtener ID del proyecto
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Debug: Ver todos los campos recibidos
	fmt.Printf("DEBUG UpdateProject - ID: %d\n", id)
	fmt.Printf("  nombreProyecto: %s\n", ctx.PostForm("nombreProyecto"))
	fmt.Printf("  fecha: %s\n", ctx.PostForm("fecha"))
	fmt.Printf("  categoria: %s\n", ctx.PostForm("categoria"))
	fmt.Printf("  descripcion: %s\n", ctx.PostForm("descripcion"))
	fmt.Printf("  lat: %s\n", ctx.PostForm("lat"))
	fmt.Printf("  lng: %s\n", ctx.PostForm("lng"))
	fmt.Printf("  userId: %s\n", ctx.PostForm("userId"))

	var project entities.Project
	project.Id = id
	project.NombreProyecto = ctx.PostForm("nombreProyecto")
	project.Fecha = ctx.PostForm("fecha")
	project.Categoria = ctx.PostForm("categoria")
	project.Descripcion = ctx.PostForm("descripcion")
	

	// Validaciones básicas
	if project.NombreProyecto == "" {
		return
	}
	if project.Categoria == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La categoría es obligatoria"})
		return
	}

	// ⚠️ CRÍTICO: Manejar UserId
	userIdStr := ctx.PostForm("userId")
	if userIdStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El userId es obligatorio"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El userId debe ser un número válido"})
		return
	}

	if userId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El userId debe ser mayor a 0"})
		return
	}

	project.UserId = userId
	fmt.Printf("DEBUG: UserId asignado correctamente: %d\n", project.UserId)

	// Coordenadas
	latStr := ctx.PostForm("lat")
	lngStr := ctx.PostForm("lng")

	if latStr != "" {
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Latitud inválida"})
			return
		}
		project.Lat = lat
	}

	if lngStr != "" {
		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Longitud inválida"})
			return
		}
		project.Lng = lng
	}

	// Manejo de imagen (opcional en update)
	var imagePath string
	file, err := ctx.FormFile("img")
	if err == nil {
		filename := "tmp_" + time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
		imagePath = filepath.Join("tmp", filename)

		if err := ctx.SaveUploadedFile(file, imagePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar imagen"})
			return
		}
		defer os.Remove(imagePath)
	}

	// Debug final antes del use case
	fmt.Printf("DEBUG: Proyecto completo antes del use case: %+v\n", project)

	// Ejecutar use case
	if err := c.useCase.Execute(project, imagePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar proyecto: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Proyecto actualizado exitosamente"})
}