package controllers

import (
	"net/http"
	"os"
	"path/filepath"
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

	// Imagen recibida
	file, err := ctx.FormFile("img")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La imagen es obligatoria"})
		return
	}

	// Guardar temporalmente la imagen
	filename := "tmp_" + time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
	path := filepath.Join("tmp", filename)

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la imagen"})
		return
	}

	// Ejecutar caso de uso con ruta de imagen
	if err := c.useCase.Execute(project, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el proyecto: " + err.Error()})
		return
	}

	// Eliminar archivo temporal
	os.Remove(path)

	ctx.JSON(http.StatusCreated, gin.H{"message": "Proyecto creado exitosamente"})
}
