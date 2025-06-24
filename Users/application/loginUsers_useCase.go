package application

import (
	"fmt"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/services"
)

type LoginUseCase struct {
	db repository.UserRepository
	jwt services.TokenManager
	bcrypt services.IBcryptService
}

func NewLoginUseCase(db repository.UserRepository, jwt services.TokenManager, bcrypt services.IBcryptService) *LoginUseCase {
	return &LoginUseCase{
		db: db,
		jwt: jwt,
		bcrypt: bcrypt,
	}
}

func (lu *LoginUseCase) Execute(email string, password string) (*entities.User,string, error) {
	user, err := lu.db.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	if !lu.bcrypt.ComparePasswords(user.Password, password) {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	token, err := lu.jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, "", fmt.Errorf("error generando token: %w", err)
	}

	return user, token, nil
}

