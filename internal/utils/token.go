package utils

import (
	"fmt"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateToken(info model.UserInfo, secret []byte, duration time.Duration) (string, error) {
	claims := model.UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func VerifyToken(tokenString string, secret []byte) (*model.UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaim)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
