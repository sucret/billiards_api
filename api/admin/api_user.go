package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type userApi struct{}

var UserApi = userApi{}

func (*userApi) List(c *gin.Context) {
	var form request.UserListReq
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	list, err := service.UserService.List(form)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}
