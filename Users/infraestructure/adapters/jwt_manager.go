package adapters

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	SecretKey string
}

func (j *JWTManager) GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) ValidateToken(token string) (bool, map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return false, nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return false, nil, nil
	}

	return true, claims, nil
}