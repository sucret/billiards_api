package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type couponApi struct{}

var CouponApi = new(couponApi)

// 优惠券列表
func (*couponApi) List(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	shopId, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		response.ValidateFail(c, "参数错误")
	}

	resp := response.CouponResp{}

	resp.CouponList, err = service.CouponService.List(0, int32(shopId))
	resp.UserCouponList, err = service.UserCouponService.List(int32(userId), 0, 0, 0)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, resp)
}
