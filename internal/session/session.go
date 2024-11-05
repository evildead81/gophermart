package session

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key")

func GenerateAuthToken(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   fmt.Sprint(userID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateAuthToken(tokenStr string) (int, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	userID, _ := strconv.Atoi(claims.Subject)
	return userID, nil
}
