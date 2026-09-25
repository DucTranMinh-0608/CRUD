package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userID int64) (string, string, time.Time, error) {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return "", "", time.Time{}, ErrMissingToken
	}

	jti := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": userID,
		"jti":     jti,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return tokenStr, jti, expiresAt, nil
}

func VerifyToken(tokenString string) (int64, string, error) {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return 0, "", ErrMissingToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Sai phương thức ký")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return 0, "", ErrInvalidToken
	}

	if !token.Valid {
		return 0, "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", ErrTokenNodata
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", ErrTokenNodata
	}

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return 0, "", ErrTokenNodata
	}

	return int64(userIDFloat), jti, nil
}
