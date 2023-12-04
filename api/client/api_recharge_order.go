package api

import (
	"billiards/pkg/config"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type rechargeOrderApi struct{}

var RechargeOrderApi = new(rechargeOrderApi)

func (*rechargeOrderApi) Price(c *gin.Context) {
	price := config.GetConfig().RechargeAmount
	response.Success(c, price)
}

func (*rechargeOrderApi) Create(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	rechargeType, _ := strconv.Atoi(c.Query("recharge_type"))

	resp, err := service.RechargeOrderService.Create(int32(userId), rechargeType)
	if err != nil {
		return
	}

	response.Success(c, resp)
}

func (*rechargeOrderApi) PayResult(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderId, _ := strconv.Atoi(c.Query("order_id"))

	status, err := service.RechargeOrderService.PayResult(int32(orderId), int32(userId))

	if err != nil {
		response.BusinessFail(c, "订单未支付")
		return
	}

	response.Success(c, status)
}
