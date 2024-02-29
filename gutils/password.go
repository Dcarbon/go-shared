package gutils

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckHashPassword(hashPassword string, password string) error
}
type stringService struct{}

func NewPasswordService() PasswordService {
	return &stringService{}
}

func (s *stringService) HashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), err
}

func (s *stringService) CheckHashPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
