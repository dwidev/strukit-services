package hash

import "golang.org/x/crypto/bcrypt"

func Password(rawPass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPass), 12)
	return string(bytes), err
}
