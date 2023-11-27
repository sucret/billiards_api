package router

import (
	api "billiards/api/external"
	"billiards/middleware"
	"github.com/gin-gonic/gin"
)

func setExternalRoute(r *gin.Engine) {
	clientRouters := r.Group("/e").Use(gin.Logger(), middleware.CustomRecovery())
	{
		// 微信支付回调
		clientRouters.POST("/wechat/pay-notify", api.WechatApi.PayNotify)
	}
}
