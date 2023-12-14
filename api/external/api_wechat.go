package api

import (
	"billiards/pkg/config"
	"billiards/pkg/qiniu"
	"billiards/pkg/tool"
	"billiards/response"
	"billiards/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
)

type wechatApi struct{}

var WechatApi = new(wechatApi)

// 支付回调
func (*wechatApi) PayNotify(c *gin.Context) {
	err := service.PaymentService.WechatPayNotify(c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (*wechatApi) Refund(c *gin.Context) {

	//form, _ := c.MultipartForm()
	//files := form.File["file"]

	header, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	tool.Dump(header)
	key := "avatar/" + tool.GenRandomString(32) + path.Ext(header.Filename)

	_ = qiniu.UploadFile(header, key)

	fmt.Println(config.GetConfig().Qiniu.Domain + key)

	//for _, v := range files {
	//	tool.Dump(v)
	//	//src, err := v.Open()
	//	//if err != nil {
	//	//	return
	//	//}
	//
	//	qiniu.GetQiniuClient(v)
	//
	//}
	//files := form.File["file"]

	//open, err := form.Open()
	//if err != nil {
	//	return
	//}

	//fmt.Println(form.Filename, 123)
	//tool.Dump(config.GetConfig().RechargeAmount)
	//wechat.GetApp("client").GenQrcode("pages/shop/detail/shopDetail", "shop_id=1")
	//if err != nil {
	//	return
	//}
	//fmt.Println(token)

	//fmt.Println(model.Time{})
	//return
	//service.PaymentService.Tt()
	//wechat.NewPayment().GetRefundDetail("202311280954078129131")
	//wechat.NewPayment().Refund()
}
