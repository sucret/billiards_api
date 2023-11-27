package api

import (
	"billiards/pkg/tool"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type userApi struct{}

var UserApi = new(userApi)

// 小程序code换session等用户信息
func (*userApi) Login(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		response.BusinessFail(c, "参数错误")
		return
	}

	if user, err := service.UserService.Login(code); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		tool.Dump(user)
		tokenData, _, err := service.JwtService.CreateToken(service.AppClientName, user)

		tool.Dump(tokenData)

		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}

		response.Success(c, tokenData)
	}
}
