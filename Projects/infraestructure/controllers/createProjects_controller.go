package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
)

// CreateProjectController
type CreateProjectController struct {
	useCase *application.CreateProjectUseCase
}

func NewCreateProjectController(useCase *application.CreateProjectUseCase) *CreateProjectController {
	return &CreateProjectController{useCase: useCase}
}

func (c *CreateProjectController) Execute(ctx *gin.Context) {
	// Debug: Ver todos los campos recibidos
	fmt.Println("DEBUG CreateProject - Campos recibidos:")
	fmt.Printf("  nombreProyecto: %s\n", ctx.PostForm("nombreProyecto"))
	fmt.Printf("  fecha: %s\n", ctx.PostForm("fecha"))
	fmt.Printf("  categoria: %s\n", ctx.PostForm("categoria"))
	fmt.Printf("  descripcion: %s\n", ctx.PostForm("descripcion"))
	fmt.Printf("  lat: %s\n", ctx.PostForm("lat"))
	fmt.Printf("  lng: %s\n", ctx.PostForm("lng"))
	fmt.Printf("  userId: %s\n", ctx.PostForm("userId"))

	var project entities.Project

	// Campos básicos
	project.NombreProyecto = ctx.PostForm("nombreProyecto")
	project.Fecha = ctx.PostForm("fecha")
	project.Categoria = ctx.PostForm("categoria")
	project.Descripcion = ctx.PostForm("descripcion")

	// Validaciones básicas
	if project.NombreProyecto == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El nombre del proyecto es obligatorio"})
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

	// Coordenadas del mapa
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

	// Manejo de imagen
	file, err := ctx.FormFile("img")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La imagen es obligatoria"})
		return
	}

	filename := "tmp_" + time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
	path := filepath.Join("tmp", filename)

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la imagen"})
		return
	}

	// Debug final antes de llamar al use case
	fmt.Printf("DEBUG: Proyecto completo antes del use case: %+v\n", project)

	// Ejecutar use case
	if err := c.useCase.Execute(project, path); err != nil {
		os.Remove(path) // Limpiar archivo temporal en caso de error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el proyecto: " + err.Error()})
		return
	}

	// Limpiar archivo temporal
	os.Remove(path)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Proyecto creado exitosamente"})
}
