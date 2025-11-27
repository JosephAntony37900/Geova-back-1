package application

import (
	"fmt"
	"regexp"
	"unicode"
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

func (lu *LoginUseCase) Execute(email string, password string) (*entities.User, string, error) {
	// 1. Validar formato del email
	if err := lu.validateEmail(email); err != nil {
		return nil, "", err
	}

	// 2. Validar requisitos de la contraseña
	if err := lu.validatePassword(password); err != nil {
		return nil, "", err
	}

	// 3. Buscar usuario por email (valida si el correo está registrado)
	user, err := lu.db.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("el correo electrónico no está registrado")
	}

	// 4. Verificar contraseña
	if !lu.bcrypt.ComparePasswords(user.Password, password) {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	// 5. Generar token
	token, err := lu.jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, "", fmt.Errorf("error generando token: %w", err)
	}

	return user, token, nil
}

// validateEmail valida el formato del correo electrónico
func (lu *LoginUseCase) validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("el correo electrónico es requerido")
	}

	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)

	if err != nil {
		return fmt.Errorf("error validando el correo electrónico")
	}

	if !matched {
		return fmt.Errorf("el formato del correo electrónico no es válido")
	}

	return nil
}

// validatePassword valida que la contraseña cumpla con los requisitos de seguridad
func (lu *LoginUseCase) validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("la contraseña es requerida")
	}

	if len(password) < 8 {
		return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
	}

	var (
		hasUpper   bool
		hasNumber  bool
		hasSpecial bool
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/"

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
		for _, special := range specialChars {
			if char == special {
				hasSpecial = true
				break
			}
		}
	}

	if !hasUpper && !hasNumber && !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula, un número y un carácter especial")
	}

	if !hasUpper && !hasNumber {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula y un número")
	}

	if !hasUpper && !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula y un carácter especial")
	}

	if !hasNumber && !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos un número y un carácter especial")
	}

	if !hasUpper {
		return fmt.Errorf("la contraseña debe contener al menos una mayúscula")
	}

	if !hasNumber {
		return fmt.Errorf("la contraseña debe contener al menos un número")
	}

	if !hasSpecial {
		return fmt.Errorf("la contraseña debe contener al menos un carácter especial (!@#$%%^&*()_+-=[]{}|;:,.<>?/)")
	}

	return nil
}