package adapters

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	cost int
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{
		cost: 12,
	}
}

func (b *Bcrypt) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (b *Bcrypt) ComparePasswords(hashedPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	return err == nil
	
}