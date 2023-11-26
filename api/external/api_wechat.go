package api

import (
	"billiards/pkg/tool"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type wechatApi struct{}

var WechatApi = new(wechatApi)

// 支付回调
func (*wechatApi) PayNotify(c *gin.Context) {
	orderNum := c.Query("order_num")
	order, err := service.OrderService.PaySuccess(orderNum)
	if err != nil {
		return
	}

	tool.Dump(order)
}
