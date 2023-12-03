package middleware

import (
	"billiards/pkg/config"
	"billiards/response"
	"billiards/service"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(c)
			c.Abort()
			return
		}
		fmt.Println("tokenStr", tokenStr)
		tokenStr = tokenStr[len(service.TokenType)+1:]

		jwtConfig := config.GetConfig().Jwt

		// Token 解析校验
		token, err := jwt.ParseWithClaims(tokenStr, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})
		if err != nil {
			response.TokenFail(c)
			c.Abort()
			return
		}

		claims := token.Claims.(*service.CustomClaims)
		// Token 发布者校验
		if claims.Issuer != GuardName {
			response.TokenFail(c)
			c.Abort()
			return
		}

		// token刷新
		// 如果token过期时间在两小时以内，则刷新token
		// 客户端拿到新的token之后替换掉原来的token
		//if claims.ExpiresAt-time.Now().Unix() < 7200 {
		//	adminId, _ := strconv.Atoi(claims.Id)
		//	admin := model.Admin{AdminID: int32(adminId)}
		//
		//	tokenData, _, err := service.JwtService.CreateToken(service.AppGuardName, admin)
		//
		//	if err != nil {
		//		response.TokenFail(c)
		//		c.Abort()
		//		return
		//	}
		//
		//	c.Header("NewToken", tokenData.AccessToken)
		//}

		c.Set("token", token)
		c.Set("userId", claims.Id)
	}
}
