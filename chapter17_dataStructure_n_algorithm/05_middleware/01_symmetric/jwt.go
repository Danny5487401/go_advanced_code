package _1_symmetric

import (
	"net/http"
	"time"

	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/05_middleware/tokenErr"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"

	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/05_middleware/models"
)

type JWTConfig struct {
	SigningKey string `json:"key"`
}

func JWTAuth(key []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		j := NewJWT(key)
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == tokenErr.TokenExpired {
				if err == tokenErr.TokenExpired {
					c.JSON(http.StatusUnauthorized, map[string]string{
						"msg": "授权已过期",
					})
					c.Abort()
					return
				}
			}

			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.UserID)
		c.Next()
	}
}

type JWT struct {
	// 在hmac中key 必须是Key must be []byte
	// 在rsa中key 必须是*rsa.PrivateKey 对象
	SigningKey []byte
}

func NewJWT(secret []byte) *JWT {
	return &JWT{
		secret, //可以设置过期时间
	}
}

// 创建一个token
func (j *JWT) CreateToken(user *models.User) (string, error) {
	claims := models.CustomClaims{
		User: user,
	}
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour)) // 过期时间，必须设置,
	claims.Issuer = "danny"
	// 1.组成token结构体 ,使用hmac方法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 2.传入 key 返回token或者error
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 带callback函数keyFunc:解析方法的回调函数 方法返回秘钥 可以根据不同的判断返回不同的秘钥
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
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
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, tokenErr.TokenInvalid

	}
	return nil, tokenErr.TokenInvalid

}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return j.CreateToken(claims.User)
	}
	return "", tokenErr.TokenInvalid
}
