package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
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
		tokenData, _, err := service.JwtService.CreateToken(service.AppClientName, user)

		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}

		resp := response.UserLogin{User: *user, TokenData: tokenData}

		response.Success(c, resp)
	}
}

func (*userApi) Save(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	form := request.SaveUser{}
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	list, err := service.UserService.Save(int32(userId), form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}
