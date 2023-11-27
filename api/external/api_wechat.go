package api

import (
	"billiards/pkg/wechat"
	"github.com/gin-gonic/gin"
)

type wechatApi struct{}

var WechatApi = new(wechatApi)

// 支付回调
func (*wechatApi) PayNotify(c *gin.Context) {

	//tool.Dump(c.Request)
	//request := notify.Request{}
	//err := c.ShouldBindJSON(&request)
	//if err != nil {
	//	//fmt.Println(err)
	//	return
	//}

	//tool.Dump(request)

	//tool.Dump(request)
	//tool.Dump(c.Request)
	wechat.NewPayment().GetPayResult(c.Request)
	//
	//orderNum := c.Query("order_num")
	//order, err := service.OrderService.PaySuccess(orderNum)
	//if err != nil {
	//	return
	//}
	//
	//tool.Dump(order)
}
