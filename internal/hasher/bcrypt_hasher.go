package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func New() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(passwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *BcryptHasher) Compare(hashedPasswd, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd))
	return err == nil
}
