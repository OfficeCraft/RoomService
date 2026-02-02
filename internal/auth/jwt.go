package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWTTokenToUserID(tokenString string, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["userId"].(string); ok {
			return userID, nil
		}
	}

	return "", errors.New("userId not found in token claims")
}
