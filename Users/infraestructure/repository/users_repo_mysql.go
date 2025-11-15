package repository

import (
	"fmt"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/core"
)

type UserMySQLRepository struct {
	db *core.Conn_MySQL
}

func NewUserMySQLRepository(db *core.Conn_MySQL) repository.UserRepository {
	return &UserMySQLRepository{
		db: db,
	}
}

// Save guarda un nuevo usuario en la base de datos
func (r *UserMySQLRepository) Save(user entities.User) error {
	// Verificar si el email ya existe
	existingUser, _ := r.FindByEmail(user.Email)
	if existingUser != nil {
		return fmt.Errorf("el email %s ya está registrado", user.Email)
	}

	query := `INSERT INTO users (Username, Nombre, Apellidos, Email, Password) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecutePreparedQuery(query,
		user.Username, user.Nombre, user.Apellidos, user.Email, user.Password)

	if err != nil {
		return fmt.Errorf("error al guardar usuario: %w", err)
	}

	return nil
}

// Update actualiza un usuario existente
func (r *UserMySQLRepository) Update(user entities.User) error {
	// Verificar que el usuario existe
	existingUser, err := r.FindById(user.Id)
	if err != nil || existingUser == nil {
		return fmt.Errorf("el usuario con ID %d no existe", user.Id)
	}

	// Verificar si el email ya está siendo usado por otro usuario
	userWithEmail, _ := r.FindByEmail(user.Email)
	if userWithEmail != nil && userWithEmail.Id != user.Id {
		return fmt.Errorf("el email %s ya está siendo usado por otro usuario", user.Email)
	}

	query := `UPDATE users SET Username = ?, Nombre = ?, Apellidos = ?, Email = ?, Password = ? WHERE Id = ?`
	_, err = r.db.ExecutePreparedQuery(query,
		user.Username, user.Nombre, user.Apellidos, user.Email, user.Password, user.Id)

	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	return nil
}

// Delete elimina un usuario por ID
func (r *UserMySQLRepository) Delete(id int) error {
	// Verificar que el usuario existe
	existingUser, err := r.FindById(id)
	if err != nil || existingUser == nil {
		return fmt.Errorf("el usuario con ID %d no existe", id)
	}

	query := `DELETE FROM users WHERE Id = ?`
	_, err = r.db.ExecutePreparedQuery(query, id)

	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	return nil
}

// FindById busca un usuario por ID
func (r *UserMySQLRepository) FindById(id int) (*entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users WHERE Id = ?`
	rows := r.db.FetchRows(query, id)
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

// FindAll obtiene todos los usuarios
func (r *UserMySQLRepository) FindAll() ([]entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users ORDER BY Id`
	rows := r.db.FetchRows(query)
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

// FindByEmail busca un usuario por email
func (r *UserMySQLRepository) FindByEmail(email string) (*entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users WHERE Email = ?`
	rows := r.db.FetchRows(query, email)
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
