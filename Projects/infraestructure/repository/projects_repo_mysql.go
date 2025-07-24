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
    // Primero verificar que el usuario existe
    userQuery := `SELECT id FROM users WHERE id = ?`
    userRows := r.db.FetchRows(userQuery, project.UserId)
    defer userRows.Close()
    
    if !userRows.Next() {
        return fmt.Errorf("el usuario con ID %d no existe", project.UserId)
    }
    
    // Si el usuario existe, proceder con la inserción
    query := `INSERT INTO projects (NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
    _, err := r.db.ExecutePreparedQuery(query, project.NombreProyecto, project.Fecha, project.Categoria, project.Descripcion, project.Img, project.Lat, project.Lng, project.UserId)
    if err != nil {
        return fmt.Errorf("error al guardar proyecto: %w", err)
    }
    return nil
}

func (r *ProjectMySQLRepository) FindById(id int) (*entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Id = ?`
	rows := r.db.FetchRows(query, id)
	defer rows.Close()

	if rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId); err != nil {
			return nil, err
		}
		return &project, nil
	}
	return nil, fmt.Errorf("proyecto no encontrado")
}

func (r *ProjectMySQLRepository) FindAll() ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects`
	rows := r.db.FetchRows(query)
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) Update(project entities.Project) error {
	query := `UPDATE projects SET NombreProyecto = ?, Fecha = ?, Categoria = ?, Descripcion = ?, Img = ?, Lat = ?, Lng = ?, user_id = ? WHERE Id = ?`
	_, err := r.db.ExecutePreparedQuery(query, project.NombreProyecto, project.Fecha, project.Categoria, project.Descripcion, project.Img, project.Lat, project.Lng, project.UserId, project.Id)
	return err
}

func (r *ProjectMySQLRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE Id = ?`
	_, err := r.db.ExecutePreparedQuery(query, id)
	return err
}

func (r *ProjectMySQLRepository) FindByName(nombre string) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE NombreProyecto LIKE ?`
	rows := r.db.FetchRows(query, "%"+nombre+"%")
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) FindByCategory(categoria string) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Categoria = ?`
	rows := r.db.FetchRows(query, categoria)
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) FindByDate(fecha string) ([]entities.Project, error) {
	fmt.Printf("DEBUG Repository - Buscando proyectos con fecha: '%s'\n", fecha)
	
	debugQuery := `SELECT DISTINCT Fecha FROM projects LIMIT 10`
	debugRows := r.db.FetchRows(debugQuery)
	fmt.Println("DEBUG - Fechas existentes en BD:")
	for debugRows.Next() {
		var existingFecha string
		debugRows.Scan(&existingFecha)
		fmt.Printf("  - '%s'\n", existingFecha)
	}
	debugRows.Close()
	
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Fecha = ?`
	rows := r.db.FetchRows(query, fecha)
	defer rows.Close()

	var projects []entities.Project
	
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(
			&project.Id, 
			&project.NombreProyecto, 
			&project.Fecha, 
			&project.Categoria, 
			&project.Descripcion, 
			&project.Img, 
			&project.Lat, 
			&project.Lng, 
			&project.UserId,
		); err != nil {
			return nil, fmt.Errorf("error escaneando fila: %w", err)
		}
		projects = append(projects, project)
		fmt.Printf("DEBUG Repository - Proyecto encontrado: %s, Fecha: %s\n", project.NombreProyecto, project.Fecha)
	}
	
	fmt.Printf("DEBUG Repository - Total proyectos encontrados con búsqueda exacta: %d\n", len(projects))
	
	// PASO 3: Si no encontró nada, intentar búsqueda flexible
	if len(projects) == 0 {
		fmt.Println("DEBUG - No se encontraron proyectos con búsqueda exacta, intentando búsqueda flexible...")
		
		// Búsqueda con LIKE para casos de espacios o caracteres extra
		flexQuery := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE TRIM(Fecha) LIKE ?`
		flexRows := r.db.FetchRows(flexQuery, "%"+fecha+"%")
		defer flexRows.Close()
		
		for flexRows.Next() {
			var project entities.Project
			if err := flexRows.Scan(
				&project.Id, 
				&project.NombreProyecto, 
				&project.Fecha, 
				&project.Categoria, 
				&project.Descripcion, 
				&project.Img, 
				&project.Lat, 
				&project.Lng, 
				&project.UserId,
			); err != nil {
				return nil, fmt.Errorf("error escaneando fila flexible: %w", err)
			}
			projects = append(projects, project)
			fmt.Printf("DEBUG Repository - Proyecto encontrado (búsqueda flexible): %s, Fecha: '%s'\n", project.NombreProyecto, project.Fecha)
		}
		
		fmt.Printf("DEBUG Repository - Total proyectos encontrados con búsqueda flexible: %d\n", len(projects))
	}
	
	return projects, nil
}

func (r *ProjectMySQLRepository) FindByUserId(userId int) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE user_id = ?`
	rows := r.db.FetchRows(query, userId)
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var project entities.Project
		if err := rows.Scan(&project.Id, &project.NombreProyecto, &project.Fecha, &project.Categoria, &project.Descripcion, &project.Img, &project.Lat, &project.Lng, &project.UserId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectMySQLRepository) SaveManyProjects(projects []entities.Project) error {
	for _, p := range projects {
		err := r.Save(p)
		if err != nil {
			return err
		}
	}
	return nil
}
