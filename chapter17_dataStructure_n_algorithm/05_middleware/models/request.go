package models

import (
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims JWT请求相关：自定义Claims
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
