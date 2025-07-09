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

type UpdateProjectController struct {
	useCase *application.UpdateProjectUseCase
}

func NewUpdateProjectController(useCase *application.UpdateProjectUseCase) *UpdateProjectController {
	return &UpdateProjectController{useCase: useCase}
}

func (c *UpdateProjectController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var project entities.Project
	project.Id = id
	project.NombreProyecto = ctx.PostForm("nombreProyecto")
	project.Fecha = ctx.PostForm("fecha")
	project.Categoria = ctx.PostForm("categoria")
	project.Descripcion = ctx.PostForm("descripcion")

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

	if err := c.useCase.Execute(project, imagePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar proyecto: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Proyecto actualizado exitosamente"})
}
