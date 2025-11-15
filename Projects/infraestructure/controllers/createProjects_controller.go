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
	
	fmt.Println("DEBUG CreateProject - Campos recibidos:")
	fmt.Printf("  nombreProyecto: %s\n", ctx.PostForm("nombreProyecto"))
	fmt.Printf("  fecha: %s\n", ctx.PostForm("fecha"))
	fmt.Printf("  categoria: %s\n", ctx.PostForm("categoria"))
	fmt.Printf("  descripcion: %s\n", ctx.PostForm("descripcion"))
	fmt.Printf("  lat: %s\n", ctx.PostForm("lat"))
	fmt.Printf("  lng: %s\n", ctx.PostForm("lng"))
	fmt.Printf("  userId: %s\n", ctx.PostForm("userId"))

	var project entities.Project

	
	project.NombreProyecto = ctx.PostForm("nombreProyecto")
	project.Fecha = ctx.PostForm("fecha")
	project.Categoria = ctx.PostForm("categoria")
	project.Descripcion = ctx.PostForm("descripcion")

	
	if project.NombreProyecto == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El nombre del proyecto es obligatorio"})
		return
	}
	if project.Categoria == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La categoría es obligatoria"})
		return
	}

	
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

	// Modo offline
	var imagePath string
	file, err := ctx.FormFile("img")
	
	if err != nil {
		
		fmt.Println("INFO: No se proporcionó imagen, creando proyecto sin imagen")
		imagePath = ""
	} else {
		
		filename := "tmp_" + time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
		imagePath = filepath.Join("tmp", filename)

		if err := ctx.SaveUploadedFile(file, imagePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la imagen temporal"})
			return
		}
		fmt.Printf("DEBUG: Imagen temporal guardada: %s\n", imagePath)
	}

	
	fmt.Printf("DEBUG: Proyecto completo antes del use case: %+v\n", project)
	fmt.Printf("DEBUG: Ruta de imagen: %s\n", imagePath)

	
	result, err := c.useCase.Execute(project, imagePath)
	
	
	if imagePath != "" {
		os.Remove(imagePath)
		fmt.Printf("DEBUG: Archivo temporal eliminado: %s\n", imagePath)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al crear el proyecto: " + err.Error(),
			"success": false,
		})
		return
	}

	
	statusCode := http.StatusCreated
	response := gin.H{
		"success":    result.Success,
		"message":    result.Message,
		"is_offline": result.IsOffline,
		"has_image":  result.HasImage,
	}

	
	if result.IsOffline {
		response["warning"] = "Proyecto creado en modo offline"
		if result.HasImage {
			response["image_note"] = "La imagen se subirá automáticamente cuando haya conexión a internet"
		}
	}

	fmt.Printf("DEBUG: Respuesta final: %+v\n", response)
	ctx.JSON(statusCode, response)
}