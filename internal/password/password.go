package password

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
