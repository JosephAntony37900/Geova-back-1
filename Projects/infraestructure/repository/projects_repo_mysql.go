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
	return &ProjectMySQLRepository{
		db: db,
	}
}

func (r *ProjectMySQLRepository) Save(project *entities.Project) error {
	query := `INSERT INTO projects (NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.ExecutePreparedQuery(query, 
		project.NombreProyecto, 
		project.Fecha, 
		project.Categoria, 
		project.Descripcion, 
		project.Img, 
		project.Lat, 
		project.Lng, 
		project.UserId)
	
	if err != nil {
		return fmt.Errorf("error al guardar proyecto: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error al obtener ID del proyecto insertado: %w", err)
	}

	project.Id = int(lastInsertID)
	
	fmt.Printf("DEBUG [MySQL]: Proyecto guardado con ID=%d\n", project.Id)
	return nil
}

func (r *ProjectMySQLRepository) Update(project *entities.Project) error {
	existingProject, err := r.FindById(project.Id)
	if err != nil || existingProject == nil {
		return fmt.Errorf("el proyecto con ID %d no existe", project.Id)
	}
	
	query := `UPDATE projects 
	          SET NombreProyecto = ?, Fecha = ?, Categoria = ?, Descripcion = ?, Img = ?, Lat = ?, Lng = ?, user_id = ? 
	          WHERE Id = ?`
	
	_, err = r.db.ExecutePreparedQuery(query, 
		project.NombreProyecto, 
		project.Fecha, 
		project.Categoria, 
		project.Descripcion, 
		project.Img, 
		project.Lat, 
		project.Lng, 
		project.UserId, 
		project.Id)
	
	if err != nil {
		return fmt.Errorf("error al actualizar proyecto: %w", err)
	}
	
	fmt.Printf("DEBUG [MySQL]: Proyecto ID=%d actualizado exitosamente\n", project.Id)
	return nil
}

func (r *ProjectMySQLRepository) Delete(id int) error {
existingProject, err := r.FindById(id)
if err != nil || existingProject == nil {
return fmt.Errorf("el proyecto con ID %d no existe", id)
}
query := `DELETE FROM projects WHERE Id = ?`
_, err = r.db.ExecutePreparedQuery(query, id)
if err != nil {
return fmt.Errorf("error al eliminar proyecto: %w", err)
}
return nil
}

func (r *ProjectMySQLRepository) FindById(id int) (*entities.Project, error) {
query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Id = ?`
rows := r.db.FetchRows(query, id)
defer rows.Close()
if rows.Next() {
var project entities.Project
err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId)
if err != nil {
return nil, err
}
return &project, nil
}
return nil, fmt.Errorf("proyecto no encontrado")
}

func (r *ProjectMySQLRepository) FindAll() ([]entities.Project, error) {
query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects ORDER BY Id DESC`
rows := r.db.FetchRows(query)
defer rows.Close()
var projects []entities.Project
for rows.Next() {
var project entities.Project
err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId)
if err != nil {
return nil, err
}
projects = append(projects, project)
}
return projects, nil
}

func (r *ProjectMySQLRepository) FindByName(nombre string) ([]entities.Project, error) {
query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE NombreProyecto LIKE ? ORDER BY Id DESC`
rows := r.db.FetchRows(query, "%"+nombre+"%")
defer rows.Close()
var projects []entities.Project
for rows.Next() {
var project entities.Project
err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId)
if err != nil {
return nil, err
}
projects = append(projects, project)
}
return projects, nil
}

func (r *ProjectMySQLRepository) FindByCategory(categoria string) ([]entities.Project, error) {
query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Categoria = ? ORDER BY Id DESC`
rows := r.db.FetchRows(query, categoria)
defer rows.Close()
var projects []entities.Project
for rows.Next() {
var project entities.Project
err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId)
if err != nil {
return nil, err
}
projects = append(projects, project)
}
return projects, nil
}

func (r *ProjectMySQLRepository) FindByDate(fecha string) ([]entities.Project, error) {
query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Fecha = ? ORDER BY Id DESC`
rows := r.db.FetchRows(query, fecha)
defer rows.Close()
var projects []entities.Project
for rows.Next() {
var project entities.Project
err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId)
if err != nil {
return nil, err
}
projects = append(projects, project)
}
return projects, nil
}

func (r *ProjectMySQLRepository) FindByUserId(userId int) ([]entities.Project, error) {
query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE user_id = ? ORDER BY Id DESC`
rows := r.db.FetchRows(query, userId)
defer rows.Close()
var projects []entities.Project
for rows.Next() {
var project entities.Project
err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId)
if err != nil {
return nil, err
}
projects = append(projects, project)
}
return projects, nil
}
