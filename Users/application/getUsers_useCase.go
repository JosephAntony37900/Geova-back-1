package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
)

type GetUsers struct {
	db repository.UserRepository
}

func NewGetUsersUseCase(db repository.UserRepository) *GetUsers {
	return &GetUsers{db: db}
}

func (gu *GetUsers) Execute() ([]entities.User, error) {
	users, err := gu.db.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}