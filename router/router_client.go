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
		clientRouter.GET("/table/get-order", api.TableApi.GetOrder)

		// 开台订单
		clientRouter.POST("/table/order/create", api.TableOrderApi.Create)
		clientRouter.GET("/table/order/pay-result", api.TableOrderApi.PayResult)
		clientRouter.GET("/table/order/list", api.TableOrderApi.List)
		clientRouter.GET("/table/order/detail", api.TableOrderApi.Detail)
		clientRouter.GET("/table/order/terminate", api.TableOrderApi.Terminate)

		// 充值订单
		clientRouter.GET("/recharge/order/create", api.RechargeOrderApi.Create)
		clientRouter.GET("/recharge/order/pay-result", api.RechargeOrderApi.PayResult)
		clientRouter.GET("/recharge/price", api.RechargeOrderApi.Price)

		// 优惠券
		clientRouter.GET("/coupon/list", api.CouponApi.List)

		// 用户优惠券
		clientRouter.GET("/user/coupon/list", api.UserCouponApi.List)

		// 统一下单接口
		clientRouter.POST("/order/create", api.OrderApi.Create)
		clientRouter.GET("/order/pay-result", api.OrderApi.PayResult)
	}
}
