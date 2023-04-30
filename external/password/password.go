package password

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Verify(hashedPwd, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
}
