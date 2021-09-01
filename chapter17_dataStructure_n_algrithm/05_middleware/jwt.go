package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"go_advanced_code/chapter17_dataStructure_n_algrithm/05_middleware/models"
)

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

var JwtConfigInfo = &JWTConfig{SigningKey: "dhsakjdhsajkdhiuw"}

func JWTAuth() gin.HandlerFunc {
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
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				if err == TokenExpired {
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
		c.Set("userId", claims.ID)
		c.Next()
	}
}

type JWT struct {
	// 在hmac中key必须是Key must be []byte
	// 在rsa中key 必须是*rsa.PrivateKey 对象
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(JwtConfigInfo.SigningKey), //可以设置过期时间
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
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
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
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
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

/*
源码分析:
	// token结构
	type Token struct {
		Raw       string                 // 保存原始token解析的时候保存
		Method    SigningMethod          // 保存签名方法 目前库里有HMAC  RSA  ECDSA
		Header    map[string]interface{} // jwt中的头部
		Claims    Claims                 // jwt中第二部分荷载，Claims是一个接口
		Signature string                 // jwt中的第三部分 签名
		Valid     bool                   // 记录token是否正确
	}

	type Claims interface {
		Valid() error
	}
	// 签名方法 所有的签名方法都会实现这个接口
	// 具体可以参考https://github.com/dgrijalva/jwt-go/blob/master/hmac.go
	type SigningMethod interface {
		// 验证token的签名，如果有限返回nil
		Verify(signingString, signature string, key interface{}) error

		// 签名方法 接受头部和荷载编码过后的字符串和签名秘钥
		// 在hmac中key必须是Key must be []byte
		// 在rsa中key 必须是*rsa.PrivateKey 对象
		Sign(signingString string, key interface{}) (string, error)

		// 返回加密方法的名字 比如'HS256'
		Alg() string
	}
	//parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(cert.PrivateKey)
*/
