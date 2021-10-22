package password

import (
	"core/sec/lib/text"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const RecommendSaltLength = 24

// HashPassword will hash the given password with salt and Pepper, then return the hashed password and salt.
func HashPassword(password string, pepper string, saltLength int) (string, string, error) {
	if password == "" {
		return "", "", fmt.Errorf("the password should not be empty")
	}
	if pepper == "" {
		return "", "", fmt.Errorf("the pepper should not be empty")
	}
	if saltLength <= 0 {
		return "", "", fmt.Errorf("the salt length should not be less than or equal to 0")
	}
	randomSalt := text.RandomString(saltLength)
	p := append([]byte(password), pepper+randomSalt...)
	bytes, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	return string(bytes), randomSalt, err
}

// CheckPasswordHash check if the given password, hashedPassword, and salt match.
func CheckPasswordHash(password, hash, pepper, salt string, passwordEncryption bool) bool {
	//fmt.Println("password: ",password)
	//fmt.Println("hash: ",hash)
	if password == "" || ((pepper == "" || salt == "") && passwordEncryption) {
		return false
	}
	if !passwordEncryption {

		return password == hash
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), append([]byte(password), pepper+salt...))
	return err == nil
}
