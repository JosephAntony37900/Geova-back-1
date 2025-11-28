package application

import (
	"fmt"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	"github.com/JosephAntony37900/Geova-back-1/Users/domain/services"
)

type LoginUseCase struct {
	db     repository.UserRepository
	jwt    services.TokenManager
	bcrypt services.IBcryptService
}

func NewLoginUseCase(db repository.UserRepository, jwt services.TokenManager, bcrypt services.IBcryptService) *LoginUseCase {
	return &LoginUseCase{
		db:     db,
		jwt:    jwt,
		bcrypt: bcrypt,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User  *entities.User
	Token string
}

func (lu *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	if input.Email == "" {
		return nil, fmt.Errorf("el correo electrónico es requerido")
	}
	if input.Password == "" {
		return nil, fmt.Errorf("la contraseña es requerida")
	}

	user, err := lu.db.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	if !lu.bcrypt.ComparePasswords(user.Password, input.Password) {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	token, err := lu.jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar sesión, intente nuevamente")
	}

	return &LoginOutput{
		User:  user,
		Token: token,
	}, nil
}