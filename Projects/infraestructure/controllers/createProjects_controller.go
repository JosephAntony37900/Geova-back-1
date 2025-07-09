package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntony37900/Geova-back-1/Projects/application"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
)

type CreateProjectController struct {
	useCase *application.CreateProjectUseCase
}

func NewCreateProjectController(useCase *application.CreateProjectUseCase) *CreateProjectController {
	return &CreateProjectController{useCase: useCase}
}

func (c *CreateProjectController) Execute(ctx *gin.Context) {
	var project entities.Project

	project.NombreProyecto = ctx.PostForm("nombreProyecto")
	project.Fecha = ctx.PostForm("fecha")
	project.Categoria = ctx.PostForm("categoria")
	project.Descripcion = ctx.PostForm("descripcion")

	// Coordenadas del mapa
	latStr := ctx.PostForm("lat")
	lngStr := ctx.PostForm("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err == nil {
		project.Lat = lat
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err == nil {
		project.Lng = lng
	}

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

	if err := c.useCase.Execute(project, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el proyecto: " + err.Error()})
		return
	}

	os.Remove(path)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Proyecto creado exitosamente"})
}
