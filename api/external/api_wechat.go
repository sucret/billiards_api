package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type wechatApi struct{}

var WechatApi = new(wechatApi)

// 支付回调
func (*wechatApi) PayNotify(c *gin.Context) {
	err := service.PaymentService.PayNotify(c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (*wechatApi) Refund(c *gin.Context) {
	service.PaymentService.Tt()
	//wechat.NewPayment().GetRefundDetail("202311280954078129131")
	//wechat.NewPayment().Refund()
}
