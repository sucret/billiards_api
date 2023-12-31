package service

import (
	"billiards/pkg/config"
	"billiards/pkg/mysql"
	"billiards/response"
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
	TokenType       = "bearer"
	AppGuardName    = "app"
	AppClientName   = "client"
	AppBusinessName = "business"
)

// CreateToken 生成 Token
func (*jwtService) CreateToken(GuardName string, user JwtUser) (tokenData response.TokenOutPut, token *jwt.Token, err error) {

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

	tokenData = response.TokenOutPut{
		AccessToken: tokenStr,
		ExpiresIn:   int(jwtConfig.JwtTtl),
		TokenType:   TokenType,
	}
	return
}
