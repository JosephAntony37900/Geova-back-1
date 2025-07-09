package repository

import (
	"fmt"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/core"
)

type ProjectMySQLRepository struct {
	db *core.Conn_MySQL
}

func NewProjectMySQLRepository(db *core.Conn_MySQL) repository.ProjectRepository {
	return &ProjectMySQLRepository{db: db}
}

func (r *ProjectMySQLRepository) Save(project entities.Project) error {
	query := `INSERT INTO projects (NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecutePreparedQuery(query, project.NombreProyecto, project.Fecha, project.Categoria, project.Descripcion, project.Img, project.Lat, project.Lng)
	if err != nil {
		return fmt.Errorf("error al guardar proyecto: %w", err)
	}
	return nil
}

func (r *ProjectMySQLRepository) FindById(id int) (*entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng FROM projects WHERE Id = ?`
	rows := r.db.FetchRows(query, id)
	defer rows.Close()

	if rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng); err != nil {
			return nil, err
		}
		return &project, nil
	}
	return nil, fmt.Errorf("proyecto no encontrado")
}

func (r *ProjectMySQLRepository) FindAll() ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng FROM projects`
	rows := r.db.FetchRows(query)
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) Update(project entities.Project) error {
	query := `UPDATE projects SET NombreProyecto = ?, Fecha = ?, Categoria = ?, Descripcion = ?, Img = ?, Lat = ?, Lng = ? WHERE Id = ?`
	_, err := r.db.ExecutePreparedQuery(query, project.NombreProyecto, project.Fecha, project.Categoria, project.Descripcion, project.Img, project.Lat, project.Lng, project.Id)
	return err
}

func (r *ProjectMySQLRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE Id = ?`
	_, err := r.db.ExecutePreparedQuery(query, id)
	return err
}

func (r *ProjectMySQLRepository) FindByName(nombre string) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng FROM projects WHERE NombreProyecto LIKE ?`
	rows := r.db.FetchRows(query, "%"+nombre+"%")
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) FindByCategory(categoria string) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng FROM projects WHERE Categoria = ?`
	rows := r.db.FetchRows(query, categoria)
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) FindByDate(fecha string) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng FROM projects WHERE Fecha = ?`
	rows := r.db.FetchRows(query, fecha)
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
