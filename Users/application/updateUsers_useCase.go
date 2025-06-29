package application

import (
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/services"
)

type UpdateUserCase struct {
	repo   repository.UserRepository
	bcrypt services.IBcryptService
}

func NewUpdateUserUseCase(repo repository.UserRepository, bcrypt services.IBcryptService) *UpdateUserCase {
	return &UpdateUserCase{repo: repo, bcrypt: bcrypt}
}

func (uc *UpdateUserCase) Execute(user entities.User) error {
	hashedPassword, err := uc.bcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return uc.repo.Update(user)
}
