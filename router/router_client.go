package router

import (
	api "billiards/api/client"
	"billiards/middleware"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

func setClientRoute(r *gin.Engine) {

	// code换session，不需要登陆
	r.GET("/c/user/login", api.UserApi.Login)

	clientRouter := r.Group("/c").
		Use(middleware.JWTAuth(service.AppClientName), gin.Logger(), middleware.CustomRecovery())
	{
		// 店铺
		clientRouter.GET("/shop/list", api.ShopApi.List)
		clientRouter.GET("/shop/detail", api.ShopApi.Detail)

		// 球桌
		clientRouter.GET("/table/detail", api.TableApi.Detail)

		// 订单
		clientRouter.POST("/order/create", api.OrderApi.Create)
		clientRouter.GET("/order/pay-result", api.OrderApi.PayResult)
		clientRouter.GET("/order/list", api.OrderApi.List)
		clientRouter.GET("/order/detail", api.OrderApi.Detail)
		clientRouter.GET("/order/terminate", api.OrderApi.Terminate)
	}
}
