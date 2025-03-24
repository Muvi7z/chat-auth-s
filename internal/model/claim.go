package model

import "github.com/golang-jwt/jwt"

const (
	ExamplePath = "/user_v1.UserV1/Get"
)

type UserClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int32  `json:"role"`
}
