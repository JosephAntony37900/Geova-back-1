package repository

import (
	"fmt"
	"net"
	"time"
	"encoding/json"
	"log"
	

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/core"
)

type ProjectMySQLRepository struct {
	localDB  *core.Conn_MySQL  
	remoteDB *core.Conn_MySQL  
}

type PendingOperation struct {
	ID        int       `json:"id"`
	Operation string    `json:"operation"` 
	ProjectID int       `json:"project_id"`
	Data      string    `json:"data"`      
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`    
}

type PendingImageOperation struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	ImagePath   string    `json:"image_path"`   
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`       
	RetryCount  int       `json:"retry_count"`
}

func NewProjectMySQLRepository(localDB *core.Conn_MySQL, remoteDB *core.Conn_MySQL) repository.ProjectRepository {
	repo := &ProjectMySQLRepository{
		localDB:  localDB,
		remoteDB: remoteDB,
	}
	
	
	repo.createPendingOperationsTable()
	repo.createPendingImageUploadsTable()
	
	
	go repo.startSyncWorker()
	
	
	go repo.initialSync()
	
	return repo
}

func (r *ProjectMySQLRepository) createPendingOperationsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS pending_sync_operations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		operation VARCHAR(10) NOT NULL,
		project_id INT,
		data TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(10) DEFAULT 'PENDING',
		retry_count INT DEFAULT 0,
		INDEX idx_status (status),
		INDEX idx_timestamp (timestamp)
	)`
	
	r.localDB.ExecutePreparedQuery(query)
}

func (r *ProjectMySQLRepository) createPendingImageUploadsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS pending_image_uploads (
		id INT AUTO_INCREMENT PRIMARY KEY,
		project_id INT NOT NULL,
		image_path VARCHAR(500) NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(20) DEFAULT 'PENDING',
		retry_count INT DEFAULT 0,
		INDEX idx_status (status),
		INDEX idx_project_id (project_id)
	)`
	
	r.localDB.ExecutePreparedQuery(query)
}

func (r *ProjectMySQLRepository) hasInternetConnection() bool {
	timeout := time.Duration(5 * time.Second)
	_, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	return err == nil
}

func (r *ProjectMySQLRepository) isRemoteDBAvailable() bool {
	if r.remoteDB == nil {
		return false
	}
	
	if !r.hasInternetConnection() {
		return false
	}
	
	err := r.remoteDB.DB.Ping()
	return err == nil
}

func (r *ProjectMySQLRepository) startSyncWorker() {
	ticker := time.NewTicker(30 * time.Second) 
	defer ticker.Stop()
	
	for range ticker.C {
		if r.isRemoteDBAvailable() {
			r.processPendingOperations()
		}
	}
}


func (r *ProjectMySQLRepository) initialSync() {
	
	time.Sleep(5 * time.Second)
	
	if !r.isRemoteDBAvailable() {
		log.Println("INFO: BD remota no disponible para sincronización inicial")
		return
	}
	
	log.Println("INFO: Iniciando sincronización inicial...")
	
	
	r.syncFromRemoteToLocal()
	
	
	r.processPendingOperations()
	
	log.Println("INFO: Sincronización inicial completada")
}


func (r *ProjectMySQLRepository) syncFromRemoteToLocal() {
	log.Println("INFO: Sincronizando desde BD remota a local...")
	
	
	remoteQuery := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects ORDER BY Id`
	remoteRows := r.remoteDB.FetchRows(remoteQuery)
	defer remoteRows.Close()
	
	syncCount := 0
	updateCount := 0
	
	for remoteRows.Next() {
		var remoteProject entities.Project
		err := remoteRows.Scan(
			&remoteProject.Id, 
			&remoteProject.NombreProyecto, 
			&remoteProject.Fecha, 
			&remoteProject.Categoria, 
			&remoteProject.Descripcion, 
			&remoteProject.Img, 
			&remoteProject.Lat, 
			&remoteProject.Lng, 
			&remoteProject.UserId,
		)
		if err != nil {
			log.Printf("ERROR: Error al escanear proyecto remoto: %v", err)
			continue
		}
		
		
		localProject, err := r.findByIdLocal(remoteProject.Id)
		
		if err != nil {
			
			if r.insertProjectToLocal(remoteProject) {
				syncCount++
				log.Printf("INFO: Proyecto %d sincronizado desde remota a local", remoteProject.Id)
			}
		} else {
			
			if r.projectNeedsUpdate(*localProject, remoteProject) {
				if r.updateProjectInLocal(remoteProject) {
					updateCount++
					log.Printf("INFO: Proyecto %d actualizado en local desde remota", remoteProject.Id)
				}
			}
		}
	}
	
	log.Printf("INFO: Sincronización completada - %d proyectos nuevos, %d actualizados", syncCount, updateCount)
}

func (r *ProjectMySQLRepository) findByIdLocal(id int) (*entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Id = ?`
	rows := r.localDB.FetchRows(query, id)
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

func (r *ProjectMySQLRepository) insertProjectToLocal(project entities.Project) bool {
	query := `INSERT INTO projects (Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.localDB.ExecutePreparedQuery(query, 
		project.Id, project.NombreProyecto, project.Fecha, project.Categoria, 
		project.Descripcion, project.Img, project.Lat, project.Lng, project.UserId)
	
	return err == nil
}

func (r *ProjectMySQLRepository) updateProjectInLocal(project entities.Project) bool {
	query := `UPDATE projects SET NombreProyecto = ?, Fecha = ?, Categoria = ?, Descripcion = ?, Img = ?, Lat = ?, Lng = ?, user_id = ? WHERE Id = ?`
	_, err := r.localDB.ExecutePreparedQuery(query, 
		project.NombreProyecto, project.Fecha, project.Categoria, 
		project.Descripcion, project.Img, project.Lat, project.Lng, 
		project.UserId, project.Id)
	
	return err == nil
}

func (r *ProjectMySQLRepository) projectNeedsUpdate(local, remote entities.Project) bool {
	return local.NombreProyecto != remote.NombreProyecto ||
		   local.Fecha != remote.Fecha ||
		   local.Categoria != remote.Categoria ||
		   local.Descripcion != remote.Descripcion ||
		   local.Img != remote.Img ||
		   local.Lat != remote.Lat ||
		   local.Lng != remote.Lng ||
		   local.UserId != remote.UserId
}


func (r *ProjectMySQLRepository) processPendingOperations() {
	query := `SELECT id, operation, project_id, data, timestamp FROM pending_sync_operations 
			  WHERE status = 'PENDING' AND retry_count < 3 
			  ORDER BY timestamp ASC LIMIT 10`
	
	rows := r.localDB.FetchRows(query)
	defer rows.Close()
	
	for rows.Next() {
		var op PendingOperation
		var data string
		var timestampStr string 
		
		err := rows.Scan(&op.ID, &op.Operation, &op.ProjectID, &data, &timestampStr)
		if err != nil {
			log.Printf("Error scanning pending operation: %v", err)
			continue
		}
		
		
		if timestampStr != "" {
			if parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", timestampStr); parseErr == nil {
				op.Timestamp = parsedTime
			} else {
				if parsedTime, parseErr := time.Parse(time.RFC3339, timestampStr); parseErr == nil {
					op.Timestamp = parsedTime
				} else {
					log.Printf("Warning: Could not parse timestamp %s: %v", timestampStr, parseErr)
					op.Timestamp = time.Now() 
				}
			}
		}
		
		success := false
		
		switch op.Operation {
		case "CREATE":
			var project entities.Project
			if json.Unmarshal([]byte(data), &project) == nil {
				success = r.saveToRemote(project)
			}
		case "UPDATE":
			var project entities.Project
			if json.Unmarshal([]byte(data), &project) == nil {
				success = r.updateInRemote(project)
			}
		case "DELETE":
			success = r.deleteFromRemote(op.ProjectID)
		}
		
		if success {
			r.markOperationAsSynced(op.ID)
			log.Printf("Successfully synced operation %d (%s)", op.ID, op.Operation)
		} else {
			r.incrementRetryCount(op.ID)
			log.Printf("Failed to sync operation %d (%s)", op.ID, op.Operation)
		}
	}
}

func (r *ProjectMySQLRepository) markOperationAsSynced(opID int) {
	query := `UPDATE pending_sync_operations SET status = 'SYNCED' WHERE id = ?`
	r.localDB.ExecutePreparedQuery(query, opID)
}

func (r *ProjectMySQLRepository) incrementRetryCount(opID int) {
	query := `UPDATE pending_sync_operations SET retry_count = retry_count + 1 WHERE id = ?`
	r.localDB.ExecutePreparedQuery(query, opID)
}

func (r *ProjectMySQLRepository) addPendingOperation(operation string, projectID int, project *entities.Project) {
	var data string
	if project != nil {
		if jsonData, err := json.Marshal(project); err == nil {
			data = string(jsonData)
		}
	}
	
	query := `INSERT INTO pending_sync_operations (operation, project_id, data) VALUES (?, ?, ?)`
	r.localDB.ExecutePreparedQuery(query, operation, projectID, data)
}






func (r *ProjectMySQLRepository) Save(project entities.Project) error {

	userQuery := `SELECT id FROM users WHERE id = ?`
	userRows := r.localDB.FetchRows(userQuery, project.UserId)
	defer userRows.Close()
	
	if !userRows.Next() {
		return fmt.Errorf("el usuario con ID %d no existe", project.UserId)
	}
	
	
	localQuery := `INSERT INTO projects (NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.localDB.ExecutePreparedQuery(localQuery, 
		project.NombreProyecto, project.Fecha, project.Categoria, 
		project.Descripcion, project.Img, project.Lat, project.Lng, project.UserId)
	
	if err != nil {
		return fmt.Errorf("error al guardar proyecto en BD local: %w", err)
	}
	
	
	if lastID, err := result.LastInsertId(); err == nil {
		project.Id = int(lastID)
	}
	
	
	if r.isRemoteDBAvailable() {
		if r.saveToRemote(project) {
			log.Printf("Proyecto %d guardado en ambas BDs exitosamente", project.Id)
		} else {
			r.addPendingOperation("CREATE", project.Id, &project)
			log.Printf("Proyecto %d guardado localmente, pendiente sincronización remota", project.Id)
		}
	} else {
		r.addPendingOperation("CREATE", project.Id, &project)
		log.Printf("Proyecto %d guardado localmente sin conexión, pendiente sincronización", project.Id)
	}
	
	return nil
}

func (r *ProjectMySQLRepository) Update(project entities.Project) error {
	
	localQuery := `UPDATE projects SET NombreProyecto = ?, Fecha = ?, Categoria = ?, Descripcion = ?, Img = ?, Lat = ?, Lng = ?, user_id = ? WHERE Id = ?`
	_, err := r.localDB.ExecutePreparedQuery(localQuery, 
		project.NombreProyecto, project.Fecha, project.Categoria, 
		project.Descripcion, project.Img, project.Lat, project.Lng, 
		project.UserId, project.Id)
	
	if err != nil {
		return fmt.Errorf("error al actualizar proyecto en BD local: %w", err)
	}
	
	
	if r.isRemoteDBAvailable() {
		if r.updateInRemote(project) {
			log.Printf("Proyecto %d actualizado en ambas BDs exitosamente", project.Id)
		} else {
			r.addPendingOperation("UPDATE", project.Id, &project)
			log.Printf("Proyecto %d actualizado localmente, pendiente sincronización remota", project.Id)
		}
	} else {
		r.addPendingOperation("UPDATE", project.Id, &project)
		log.Printf("Proyecto %d actualizado localmente sin conexión, pendiente sincronización", project.Id)
	}
	
	return nil
}

func (r *ProjectMySQLRepository) Delete(id int) error {
	
	localQuery := `DELETE FROM projects WHERE Id = ?`
	_, err := r.localDB.ExecutePreparedQuery(localQuery, id)
	
	if err != nil {
		return fmt.Errorf("error al eliminar proyecto de BD local: %w", err)
	}
	
	
	if r.isRemoteDBAvailable() {
		if r.deleteFromRemote(id) {
			log.Printf("Proyecto %d eliminado de ambas BDs exitosamente", id)
		} else {
			r.addPendingOperation("DELETE", id, nil)
			log.Printf("Proyecto %d eliminado localmente, pendiente sincronización remota", id)
		}
	} else {
		r.addPendingOperation("DELETE", id, nil)
		log.Printf("Proyecto %d eliminado localmente sin conexión, pendiente sincronización", id)
	}
	
	return nil
}


func (r *ProjectMySQLRepository) FindById(id int) (*entities.Project, error) {
	return r.findByIdLocal(id)
}

func (r *ProjectMySQLRepository) FindAll() ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects ORDER BY Id`
	rows := r.localDB.FetchRows(query)
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

func (r *ProjectMySQLRepository) FindByName(nombre string) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE NombreProyecto LIKE ?`
	rows := r.localDB.FetchRows(query, "%"+nombre+"%")
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
	rows := r.localDB.FetchRows(query, categoria)
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
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE Fecha = ?`
	rows := r.localDB.FetchRows(query, fecha)
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

func (r *ProjectMySQLRepository) FindByUserId(userId int) ([]entities.Project, error) {
	query := `SELECT Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id FROM projects WHERE user_id = ?`
	rows := r.localDB.FetchRows(query, userId)
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


func (r *ProjectMySQLRepository) saveToRemote(project entities.Project) bool {
	query := `INSERT INTO projects (Id, NombreProyecto, Fecha, Categoria, Descripcion, Img, Lat, Lng, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) 
			  ON DUPLICATE KEY UPDATE 
			  NombreProyecto=VALUES(NombreProyecto), Fecha=VALUES(Fecha), Categoria=VALUES(Categoria), 
			  Descripcion=VALUES(Descripcion), Img=VALUES(Img), Lat=VALUES(Lat), Lng=VALUES(Lng), user_id=VALUES(user_id)`
	
	_, err := r.remoteDB.ExecutePreparedQuery(query, project.Id, project.NombreProyecto, project.Fecha, 
		project.Categoria, project.Descripcion, project.Img, project.Lat, project.Lng, project.UserId)
	
	return err == nil
}

func (r *ProjectMySQLRepository) updateInRemote(project entities.Project) bool {
	query := `UPDATE projects SET NombreProyecto = ?, Fecha = ?, Categoria = ?, Descripcion = ?, Img = ?, Lat = ?, Lng = ?, user_id = ? WHERE Id = ?`
	_, err := r.remoteDB.ExecutePreparedQuery(query, project.NombreProyecto, project.Fecha, project.Categoria, 
		project.Descripcion, project.Img, project.Lat, project.Lng, project.UserId, project.Id)
	
	return err == nil
}

func (r *ProjectMySQLRepository) deleteFromRemote(id int) bool {
	query := `DELETE FROM projects WHERE Id = ?`
	_, err := r.remoteDB.ExecutePreparedQuery(query, id)
	return err == nil
}


