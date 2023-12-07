package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type couponApi struct{}

var CouponApi = new(couponApi)

// 优惠券列表
func (*couponApi) List(c *gin.Context) {
	list, err := service.CouponService.List(0, 0)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}
