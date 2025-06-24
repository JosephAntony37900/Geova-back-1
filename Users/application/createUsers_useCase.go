package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
)

type CreateUserUseCase struct {
	db repository.UserRepository
}

func NewCreateUserUseCase(db repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{db: db}
}

func (uc *CreateUserUseCase) Execute(user entities.User) error {
	return uc.db.Save(user)
}