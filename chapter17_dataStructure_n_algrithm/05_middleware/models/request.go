package models

import (
	"github.com/dgrijalva/jwt-go"
)

// JWT请求相关
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

//  valid（）方法
/*
	// Validate Claims
	if !p.SkipClaimsValidation {
		//调用
		if err := token.Claims.Valid(); err != nil {

			// If the Claims Valid returned an error, check if it is a validation error,
			// If it was another error type, create a ValidationError with a generic ClaimsInvalid flag set
			if e, ok := err.(*ValidationError); !ok {
				vErr = &ValidationError{Inner: err, Errors: ValidationErrorClaimsInvalid}
			} else {
				vErr = e
			}
		}
	}
*/
