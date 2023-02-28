package util

import (
	"douyin/web/initialize/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

const (
	invalidToken = "invalid token"
)

var mySecret = []byte(config.Config.Jwt.Secret)

//额外自定义使用两个字段

type MyClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New(invalidToken)
}
