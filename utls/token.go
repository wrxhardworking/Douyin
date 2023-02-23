package utls

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JwtCustomClaims 定义需求 自定义+jwt预定义
type JwtCustomClaims struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"admin"`
	jwt.StandardClaims
}

// 密钥
var jwtKey = []byte("secret")

// GenerateToken  生成token
func GenerateToken(username string, userid int64) (string, error) {

	// Set custom claims
	claims := &JwtCustomClaims{
		userid,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	//生成token并且指定私钥
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return token, err
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *JwtCustomClaims, error) {
	claims := &JwtCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
