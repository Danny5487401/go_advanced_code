package _2_asymmetric

import (
	"crypto/rsa"
	"github.com/Danny5487401/go_advanced_code/chapter12_net/08_middleware/models"
	"log"
	"os"
	"time"

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
	signBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := os.ReadFile(pubKeyPath)
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
		return nil, err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
