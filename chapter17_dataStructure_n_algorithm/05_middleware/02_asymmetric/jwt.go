package _2_asymmetric

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/05_middleware/models"
	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/05_middleware/tokenErr"
	"github.com/golang-jwt/jwt/v4"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func GenerateJWT(user *models.User) (tokenString string, err error) {

	claims := models.CustomClaims{
		User: user,
	}
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour)) // 过期时间，必须设置
	claims.Issuer = "danny"

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)

}

func InitJWT(privateKeyPath, pubKeyPath string) {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, tokenErr.TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, tokenErr.TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, tokenErr.TokenNotValidYet
			} else {
				return nil, tokenErr.TokenInvalid
			}
		}
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, tokenErr.TokenInvalid

	}
	return nil, tokenErr.TokenInvalid
}
