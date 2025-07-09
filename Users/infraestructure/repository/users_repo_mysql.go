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
	return &UserMySQLRepository{db: db}
}

func (r *UserMySQLRepository) Save(user entities.User) error {
	query := `INSERT INTO users (Username, Nombre, Apellidos, Email, Password) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecutePreparedQuery(query,
		user.Username,
		user.Nombre,
		user.Apellidos,
		user.Email,
		user.Password,
	)
	if err != nil {
		return fmt.Errorf("error al guardar usuario: %w", err)
	}
	return nil
}

func (r *UserMySQLRepository) FindById(id int) (*entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users WHERE Id = ?`
	rows := r.db.FetchRows(query, id)
	defer rows.Close()

	if rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Nombre,
			&user.Apellidos,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, fmt.Errorf("usuario no encontrado")
}

func (r *UserMySQLRepository) FindAll() ([]entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users`
	rows := r.db.FetchRows(query)
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Nombre,
			&user.Apellidos,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserMySQLRepository) FindByEmail(email string) (*entities.User, error) {
	query := `SELECT Id, Username, Nombre, Apellidos, Email, Password FROM users WHERE Email = ?`
	rows := r.db.FetchRows(query, email)
	defer rows.Close()

	if rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Nombre,
			&user.Apellidos,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, fmt.Errorf("usuario no encontrado")
}

func (r *UserMySQLRepository) Update(user entities.User) error {
	query := `UPDATE users SET Username = ?, Nombre = ?, Apellidos = ?, Email = ?, Password = ? WHERE Id = ?`
	_, err := r.db.ExecutePreparedQuery(query,
		user.Username,
		user.Nombre,
		user.Apellidos,
		user.Email,
		user.Password,
		user.Id,
	)
	return err
}

func (r *UserMySQLRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE Id = ?`
	_, err := r.db.ExecutePreparedQuery(query, id)
	return err
}
