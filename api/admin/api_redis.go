package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type redisApi struct{}

var RedisApi = new(redisApi)

func (*redisApi) Query(c *gin.Context) {
	var form request.RedisQuery
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	val, err := service.RedisService.Query(form.Method, form.Query)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, val)
}
