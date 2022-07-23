package models

import (
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims JWT请求相关：自定义Claims
type CustomClaims struct {
	*User
	jwt.RegisteredClaims
}

type User struct {
	UserID   uint64
	NickName string
}
