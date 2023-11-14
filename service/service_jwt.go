package service

import (
	"gin-api/pkg/config"
	"gin-api/pkg/mysql"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type jwtService struct {
	db *gorm.DB
}

var JwtService = &jwtService{
	db: mysql.GetDB(),
}

// JwtUser 所有需要颁发 token 的用户模型必须实现这个接口
type JwtUser interface {
	GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
	jwt.StandardClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// CreateToken 生成 Token
func (*jwtService) CreateToken(GuardName string, user JwtUser) (tokenData TokenOutPut, token *jwt.Token, err error) {

	jwtConfig := config.GetConfig().Jwt

	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + jwtConfig.JwtTtl,
				Id:        user.GetUid(),
				Issuer:    GuardName, // 用于在中间件中区分不同客户端颁发的 token，避免 token 跨端使用
				NotBefore: time.Now().Unix() - 1000,
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(jwtConfig.Secret))

	tokenData = TokenOutPut{
		tokenStr,
		int(jwtConfig.JwtTtl),
		TokenType,
	}
	return
}
