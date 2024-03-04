package util

import (
	"chsir-zy/msg-push-center/config"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// const SECRET = "123456" //认证的密码 TODO

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

		return []byte(config.JWT_KEY), nil
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

// 生成token
func GenToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": "20942",
	})
	tokenStr, err := token.SignedString([]byte(config.JWT_KEY))

	if err != nil {
		fmt.Println(err)
	}

	return tokenStr
}
