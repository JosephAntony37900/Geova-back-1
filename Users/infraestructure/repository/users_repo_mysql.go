package repository

import (
	"fmt"
	"net"
	"time"
	"encoding/json"
	"log"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/core"
)

type UserMySQLRepository struct {
	localDB  *core.Conn_MySQL  
	remoteDB *core.Conn_MySQL  
}

type PendingUserOperation struct {
	ID        int       `json:"id"`
	Operation string    `json:"operation"` // CREATE, UPDATE, DELETE
	UserID    int       `json:"user_id"`
	Data      string    `json:"data"`      // JSON serializado del usuario
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`    // PENDING, SYNCED, FAILED
}

func NewUserMySQLRepository(localDB *core.Conn_MySQL, remoteDB *core.Conn_MySQL) repository.UserRepository {
	repo := &UserMySQLRepository{
		localDB:  localDB,
		remoteDB: remoteDB,
	}
	
	// Crear tabla de operaciones pendientes
	repo.createPendingUserOperationsTable()
	
	// Iniciar worker de sincronización en background
	go repo.startUserSyncWorker()
	
	// Ejecutar sincronización inicial
	go repo.initialUserSync()
	
	return repo
}

func (r *UserMySQLRepository) createPendingUserOperationsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS pending_user_sync_operations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		operation VARCHAR(10) NOT NULL,
		user_id INT,
		data TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(10) DEFAULT 'PENDING',
		retry_count INT DEFAULT 0,
		INDEX idx_user_status (status),
		INDEX idx_user_timestamp (timestamp)
	)`
	
	r.localDB.ExecutePreparedQuery(query)
}

func (r *UserMySQLRepository) hasInternetConnection() bool {
	timeout := time.Duration(5 * time.Second)
	_, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	return err == nil
}

func (r *UserMySQLRepository) isRemoteDBAvailable() bool {
	if r.remoteDB == nil {
		return false
	}
	
	if !r.hasInternetConnection() {
		return false
	}
	
	err := r.remoteDB.DB.Ping()
	return err == nil
}

func (r *UserMySQLRepository) startUserSyncWorker() {
	ticker := time.NewTicker(10 * time.Second) 
	defer ticker.Stop()
	
	syncFromRemoteCounter := 0
	
	for range ticker.C {
		if r.isRemoteDBAvailable() {
			// Procesar operaciones pendientes (local -> remoto)
			r.processPendingUserOperations()
			
			// Cada 2 ciclos (60 segundos), sincronizar desde remoto (remoto -> local)
			syncFromRemoteCounter++
			if syncFromRemoteCounter >= 2 {
				r.syncUsersFromRemoteToLocal()
				syncFromRemoteCounter = 0
			}
		}
	}
}

//ok

func (r *UserMySQLRepository) initialUserSync() {
	// Esperar un poco antes de iniciar la sincronización
	time.Sleep(5 * time.Second)
	
	if !r.isRemoteDBAvailable() {
		log.Println("INFO: BD remota no disponible para sincronización inicial de usuarios")
		return
	}
	
	log.Println("INFO: Iniciando sincronización inicial de usuarios...")
	
	// Sincronizar desde remoto a local
	r.syncUsersFromRemoteToLocal()
	
	// Procesar operaciones pendientes
	r.processPendingUserOperations()
	
	log.Println("INFO: Sincronización inicial de usuarios completada")
}

func (r *UserMySQLRepository) syncUsersFromRemoteToLocal() {
	log.Println("INFO: Sincronizando usuarios desde BD remota a local...")
	
	// Obtener todos los usuarios de la BD remota
	remoteQuery := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users ORDER BY Id`
	remoteRows := r.remoteDB.FetchRows(remoteQuery)
	defer remoteRows.Close()
	
	syncCount := 0
	updateCount := 0
	
	for remoteRows.Next() {
		var remoteUser entities.User
		err := remoteRows.Scan(
			&remoteUser.Id, 
			&remoteUser.Username, 
			&remoteUser.Nombre, 
			&remoteUser.Apellidos, 
			&remoteUser.Email, 
			&remoteUser.Password,
		)
		if err != nil {
			log.Printf("ERROR: Error al escanear usuario remoto: %v", err)
			continue
		}
		
		// Verificar si el usuario existe en local
		localUser, err := r.findByIdLocal(remoteUser.Id)
		
		if err != nil {
			// No existe en local, insertarlo
			if r.insertUserToLocal(remoteUser) {
				syncCount++
				log.Printf("INFO: Usuario %d sincronizado desde remota a local", remoteUser.Id)
			}
		} else {
			// Existe en local, verificar si necesita actualización
			if r.userNeedsUpdate(*localUser, remoteUser) {
				if r.updateUserInLocal(remoteUser) {
					updateCount++
					log.Printf("INFO: Usuario %d actualizado en local desde remota", remoteUser.Id)
				}
			}
		}
	}
	
	log.Printf("INFO: Sincronización de usuarios completada - %d usuarios nuevos, %d actualizados", syncCount, updateCount)
}

func (r *UserMySQLRepository) findByIdLocal(id int) (*entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users WHERE Id = ?`
	rows := r.localDB.FetchRows(query, id)
	defer rows.Close()

	if rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Nombre, &user.Apellidos, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, fmt.Errorf("usuario no encontrado")
}

func (r *UserMySQLRepository) insertUserToLocal(user entities.User) bool {
	query := `INSERT INTO users (Id, Username, Nombre, Apellidos, Email, Password) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.localDB.ExecutePreparedQuery(query, 
		user.Id, user.Username, user.Nombre, user.Apellidos, user.Email, user.Password)
	
	return err == nil
}

func (r *UserMySQLRepository) updateUserInLocal(user entities.User) bool {
	query := `UPDATE users SET Username = ?, Nombre = ?, Apellidos = ?, Email = ?, Password = ? WHERE Id = ?`
	_, err := r.localDB.ExecutePreparedQuery(query, 
		user.Username, user.Nombre, user.Apellidos, user.Email, user.Password, user.Id)
	
	return err == nil
}

func (r *UserMySQLRepository) userNeedsUpdate(local, remote entities.User) bool {
	return local.Username != remote.Username ||
		   local.Nombre != remote.Nombre ||
		   local.Apellidos != remote.Apellidos ||
		   local.Email != remote.Email ||
		   local.Password != remote.Password
}

func (r *UserMySQLRepository) processPendingUserOperations() {
	query := `SELECT id, operation, user_id, data, timestamp FROM pending_user_sync_operations 
			  WHERE status = 'PENDING' AND retry_count < 3 
			  ORDER BY timestamp ASC LIMIT 10`
	
	rows := r.localDB.FetchRows(query)
	defer rows.Close()
	
	for rows.Next() {
		var op PendingUserOperation
		var data string
		var timestampStr string 
		
		err := rows.Scan(&op.ID, &op.Operation, &op.UserID, &data, &timestampStr)
		if err != nil {
			log.Printf("Error scanning pending user operation: %v", err)
			continue
		}
		
		// Parsear timestamp
		if timestampStr != "" {
			if parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", timestampStr); parseErr == nil {
				op.Timestamp = parsedTime
			} else {
				if parsedTime, parseErr := time.Parse(time.RFC3339, timestampStr); parseErr == nil {
					op.Timestamp = parsedTime
				} else {
					log.Printf("Warning: Could not parse user operation timestamp %s: %v", timestampStr, parseErr)
					op.Timestamp = time.Now() 
				}
			}
		}
		
		success := false
		
		switch op.Operation {
		case "CREATE":
			var user entities.User
			if json.Unmarshal([]byte(data), &user) == nil {
				success = r.saveToRemote(user)
			}
		case "UPDATE":
			var user entities.User
			if json.Unmarshal([]byte(data), &user) == nil {
				success = r.updateInRemote(user)
			}
		case "DELETE":
			success = r.deleteFromRemote(op.UserID)
		}
		
		if success {
			r.markUserOperationAsSynced(op.ID)
			log.Printf("Successfully synced user operation %d (%s)", op.ID, op.Operation)
		} else {
			r.incrementUserRetryCount(op.ID)
			log.Printf("Failed to sync user operation %d (%s)", op.ID, op.Operation)
		}
	}
}

func (r *UserMySQLRepository) markUserOperationAsSynced(opID int) {
	query := `UPDATE pending_user_sync_operations SET status = 'SYNCED' WHERE id = ?`
	r.localDB.ExecutePreparedQuery(query, opID)
}

func (r *UserMySQLRepository) incrementUserRetryCount(opID int) {
	query := `UPDATE pending_user_sync_operations SET retry_count = retry_count + 1 WHERE id = ?`
	r.localDB.ExecutePreparedQuery(query, opID)
}

func (r *UserMySQLRepository) addPendingUserOperation(operation string, userID int, user *entities.User) {
	var data string
	if user != nil {
		if jsonData, err := json.Marshal(user); err == nil {
			data = string(jsonData)
		}
	}
	
	query := `INSERT INTO pending_user_sync_operations (operation, user_id, data) VALUES (?, ?, ?)`
	r.localDB.ExecutePreparedQuery(query, operation, userID, data)
}

// Implementación de métodos de la interfaz UserRepository

func (r *UserMySQLRepository) Save(user entities.User) error {
	// Verificar si el email ya existe (para evitar duplicados)
	existingUser, _ := r.FindByEmail(user.Email)
	if existingUser != nil {
		return fmt.Errorf("el email %s ya está registrado", user.Email)
	}
	
	// Guardar en BD local
	localQuery := `INSERT INTO users (Username, Nombre, Apellidos, Email, Password) VALUES (?, ?, ?, ?, ?)`
	result, err := r.localDB.ExecutePreparedQuery(localQuery, 
		user.Username, user.Nombre, user.Apellidos, user.Email, user.Password)
	
	if err != nil {
		return fmt.Errorf("error al guardar usuario en BD local: %w", err)
	}
	
	// Obtener el ID generado
	if lastID, err := result.LastInsertId(); err == nil {
		user.Id = int(lastID)
	}
	
	// Intentar sincronizar con BD remota
	if r.isRemoteDBAvailable() {
		if r.saveToRemote(user) {
			log.Printf("Usuario %d guardado en ambas BDs exitosamente", user.Id)
		} else {
			r.addPendingUserOperation("CREATE", user.Id, &user)
			log.Printf("Usuario %d guardado localmente, pendiente sincronización remota", user.Id)
		}
	} else {
		r.addPendingUserOperation("CREATE", user.Id, &user)
		log.Printf("Usuario %d guardado localmente sin conexión, pendiente sincronización", user.Id)
	}
	
	return nil
}

func (r *UserMySQLRepository) Update(user entities.User) error {
	// Verificar que el usuario existe
	existingUser, err := r.FindById(user.Id)
	if err != nil || existingUser == nil {
		return fmt.Errorf("el usuario con ID %d no existe", user.Id)
	}
	
	// Verificar que el email no esté siendo usado por otro usuario
	userWithEmail, _ := r.FindByEmail(user.Email)
	if userWithEmail != nil && userWithEmail.Id != user.Id {
		return fmt.Errorf("el email %s ya está siendo usado por otro usuario", user.Email)
	}
	
	// Actualizar en BD local
	localQuery := `UPDATE users SET Username = ?, Nombre = ?, Apellidos = ?, Email = ?, Password = ? WHERE Id = ?`
	_, err = r.localDB.ExecutePreparedQuery(localQuery, 
		user.Username, user.Nombre, user.Apellidos, user.Email, user.Password, user.Id)
	
	if err != nil {
		return fmt.Errorf("error al actualizar usuario en BD local: %w", err)
	}
	
	// Intentar sincronizar con BD remota
	if r.isRemoteDBAvailable() {
		if r.updateInRemote(user) {
			log.Printf("Usuario %d actualizado en ambas BDs exitosamente", user.Id)
		} else {
			r.addPendingUserOperation("UPDATE", user.Id, &user)
			log.Printf("Usuario %d actualizado localmente, pendiente sincronización remota", user.Id)
		}
	} else {
		r.addPendingUserOperation("UPDATE", user.Id, &user)
		log.Printf("Usuario %d actualizado localmente sin conexión, pendiente sincronización", user.Id)
	}
	
	return nil
}

func (r *UserMySQLRepository) Delete(id int) error {
	// Verificar que el usuario existe
	existingUser, err := r.FindById(id)
	if err != nil || existingUser == nil {
		return fmt.Errorf("el usuario con ID %d no existe", id)
	}
	
	// Eliminar de BD local
	localQuery := `DELETE FROM users WHERE Id = ?`
	_, err = r.localDB.ExecutePreparedQuery(localQuery, id)
	
	if err != nil {
		return fmt.Errorf("error al eliminar usuario de BD local: %w", err)
	}
	
	// Intentar sincronizar con BD remota
	if r.isRemoteDBAvailable() {
		if r.deleteFromRemote(id) {
			log.Printf("Usuario %d eliminado de ambas BDs exitosamente", id)
		} else {
			r.addPendingUserOperation("DELETE", id, nil)
			log.Printf("Usuario %d eliminado localmente, pendiente sincronización remota", id)
		}
	} else {
		r.addPendingUserOperation("DELETE", id, nil)
		log.Printf("Usuario %d eliminado localmente sin conexión, pendiente sincronización", id)
	}
	
	return nil
}

func (r *UserMySQLRepository) FindById(id int) (*entities.User, error) {
	return r.findByIdLocal(id)
}

func (r *UserMySQLRepository) FindAll() ([]entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users ORDER BY Id`
	rows := r.localDB.FetchRows(query)
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Nombre, &user.Apellidos, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserMySQLRepository) FindByEmail(email string) (*entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users WHERE Email = ?`
	rows := r.localDB.FetchRows(query, email)
	defer rows.Close()

	if rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Nombre, &user.Apellidos, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, fmt.Errorf("usuario no encontrado")
}

// Métodos privados para operaciones remotas

func (r *UserMySQLRepository) saveToRemote(user entities.User) bool {
	query := `INSERT INTO users (Id, Username, Nombre, Apellidos, Email, Password) VALUES (?, ?, ?, ?, ?, ?) 
			  ON DUPLICATE KEY UPDATE 
			  Username=VALUES(Username), Nombre=VALUES(Nombre), Apellidos=VALUES(Apellidos), 
			  Email=VALUES(Email), Password=VALUES(Password)`
	
	_, err := r.remoteDB.ExecutePreparedQuery(query, user.Id, user.Username, user.Nombre, 
		user.Apellidos, user.Email, user.Password)
	
	return err == nil
}

func (r *UserMySQLRepository) updateInRemote(user entities.User) bool {
	query := `UPDATE users SET Username = ?, Nombre = ?, Apellidos = ?, Email = ?, Password = ? WHERE Id = ?`
	_, err := r.remoteDB.ExecutePreparedQuery(query, user.Username, user.Nombre, user.Apellidos, 
		user.Email, user.Password, user.Id)
	
	return err == nil
}

func (r *UserMySQLRepository) deleteFromRemote(id int) bool {
	query := `DELETE FROM users WHERE Id = ?`
	_, err := r.remoteDB.ExecutePreparedQuery(query, id)
	return err == nil
}