package util

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const SECRET = "123456" //认证的密码 TODO

type JWTAuthenticator struct{}

type Authenticator interface {
	Authenticate(c *gin.Context) (string, error)
}

var _ Authenticator = &JWTAuthenticator{}

// 认证相关的
func (jwtauth *JWTAuthenticator) Authenticate(c *gin.Context) (string, error) {
	// 获取token
	tokenStr := c.Query("token")
	if tokenStr == "" {
		return "", errors.New("token can not be empty")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(SECRET), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["uid"].(string), nil
	} else {
		return "", err
	}
}
