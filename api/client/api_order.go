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

func (*orderApi) Terminate(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	order, err := service.OrderService.Terminate(userId, orderId)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (*orderApi) Detail(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	order, err := service.OrderService.Detail(userId, orderId)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (*orderApi) List(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderType, err := strconv.Atoi(c.Query("type"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	list, err := service.OrderService.List(userId, orderType)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}

func (*orderApi) Create(c *gin.Context) {
	var form request.OrderCreate
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	userId, _ := strconv.Atoi(c.GetString("userId"))

	order, err := service.OrderService.Create(form.TableID, int32(userId))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, order)
}

// 查询用户当前订单的支付结果
func (*orderApi) PayResult(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))

	orderNum := c.Query("order_num")
	if orderNum == "" {
		response.BusinessFail(c, "参数错误")
		return
	}

	status, err := service.OrderService.GetPayStatus(orderNum, int32(userId))
	if err != nil {
		response.BusinessFail(c, "订单未支付")
		return
	}

	response.Success(c, status)
}
