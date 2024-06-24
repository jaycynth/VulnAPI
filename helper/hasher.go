package helper

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type BcryptHasher struct{}

func (b BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

type MockHasher struct{}

func (m MockHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return []byte("$2a$10$e0MYzXyjpJS7Pd0RVvHwHeFUpH0lHLjZ0f72mMBp/Rs8C3W.rFYTO"), nil
}
