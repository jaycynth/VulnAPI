package helper

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	SecretKey = []byte(os.Getenv("JWT_TOKEN"))
)

type Claims struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
	jwt.StandardClaims
}

func GenerateToken(username string, userId string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		UserId:   userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
