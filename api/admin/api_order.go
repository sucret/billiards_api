package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type orderApi struct{}

var OrderApi = new(orderApi)

func (*orderApi) List(c *gin.Context) {
	var form request.OrderList
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	list, err := service.OrderService.List(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}

func (*orderApi) Detail(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	detail, err := service.OrderService.Detail(int32(orderId), -1)
	if err != nil {
		return
	}

	response.Success(c, detail)
}
