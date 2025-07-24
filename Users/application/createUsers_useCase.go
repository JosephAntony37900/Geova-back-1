package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/services"
)

type CreateUserUseCase struct {
	db      repository.UserRepository
	bcrypt  services.IBcryptService
}

func NewCreateUserUseCase(db repository.UserRepository, bcrypt services.IBcryptService) *CreateUserUseCase {
	return &CreateUserUseCase{
		db:     db,
		bcrypt: bcrypt,
	}
}

func (uc *CreateUserUseCase) Execute(user entities.User) error {
	hashedPassword, err := uc.bcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return uc.db.Save(user)
}
