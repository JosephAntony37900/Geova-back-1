package application

import (

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
)

type GetUserById struct {
	db repository.UserRepository
}

func NewGetUserByIdUseCase(db repository.UserRepository) *GetUserById{
	return &GetUserById{db:db}
}

func (gubi *GetUserById) Execute(id int) (*entities.User, error) {
	user, err := gubi.db.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
