package util

import (
	"douyin/service/user/initialize/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var mySecret = []byte(config.Config.Jwt.Secret)

//额外自定义使用两个字段

type MyClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(username string, userID int64) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		username,
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2400).Unix(), // 过期时间
			Issuer:    config.Config.Jwt.Issuer,                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}
