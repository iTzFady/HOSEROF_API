package services

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey []byte

func init() {
	_ = godotenv.Load()
	SecretKey = []byte(os.Getenv("JWT_KEY"))
	if len(SecretKey) == 0 {
		log.Fatal("JWT_KEY is missing from environment variables")
	}
}

func jwtGenerator(id, class, role, name string) (string, error) {

	claims := jwt.MapClaims{
		"user_ID":    id,
		"user_class": class,
		"user_name":  name,
		"role":       role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(SecretKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return signed, nil
}
