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

func (*orderApi) Create(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	var form request.OrderCreate
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	resp, err := service.OrderService.Create(int32(userId), form)
	if err != nil {
		return
	}

	response.Success(c, resp)
}

func (*orderApi) PayResult(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, _ := strconv.Atoi(c.Query("order_id"))

	resp, err := service.OrderService.PayResult(int32(orderId), int32(userId))
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// 终止订单
func (*orderApi) Terminate(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		response.ValidateFail(c, "参数错误")
	}

	err = service.OrderService.Terminate(int32(orderId), int32(userId))
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (*orderApi) Detail(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		response.ValidateFail(c, "参数错误")
	}

	detail, _ := service.OrderService.Detail(int32(orderId), int32(userId))

	response.Success(c, detail)
}

// 订单续费
func (*orderApi) Renewal(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		response.ValidateFail(c, "参数错误")
	}

	resp, err := service.OrderService.Renewal(int32(orderId), int32(userId))
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	response.Success(c, resp)
}
