package util

import (
	"golang.org/x/crypto/bcrypt"
)

func init() {
	type Config struct {
		Pepper     string `json:"pepper"`
		SaltLength int    `json:"salt_length"`
	}
	var config Config
	err := ParseFileJsonConfig(&config, "config/password.config.json")
	if err != nil {
		panic(err.Error())
	}
	pepper = config.Pepper
	randomSaltLength = config.SaltLength
}

var pepper string
var randomSaltLength int

// HashPassword will hash the given password with salt and pepper, then return the hashed password and salt.
func HashPassword(password string) (string, string, error) {
	randomSalt := RandomString(randomSaltLength)
	p := append([]byte(password), pepper+randomSalt...)
	bytes, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	return string(bytes), randomSalt, err
}

// CheckPasswordHash check if the given password, hashedPassword, and salt match.
func CheckPasswordHash(password, hash, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), append([]byte(password), pepper+salt...))
	return err == nil
}
