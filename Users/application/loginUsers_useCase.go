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
		return nil, fmt.Errorf("El coorreo electrónico es requerido")
	}
	if input.Password == "" {
		return nil, fmt.Errorf("La contraseña es requerida")
	}

	user, err := lu.db.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("Correo electrónico no registrado")
	}

	if !lu.bcrypt.ComparePasswords(user.Password, input.Password) {
		return nil, fmt.Errorf("Contraseña inválida")
	}

	token, err := lu.jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("Error al iniciar sesión, intente nuevamente")
	}

	return &LoginOutput{
		User:  user,
		Token: token,
	}, nil
}