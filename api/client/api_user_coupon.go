package api

import (
	"billiards/pkg/mysql/model"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type userCouponApi struct{}

var UserCouponApi = new(userCouponApi)

// 用户优惠券列表
func (*userCouponApi) List(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))

	list, err := service.UserCouponService.List(int32(userId), 0, 0, 0)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}

// 通过店铺查询
func (*userCouponApi) GetByShop(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	shopId, err := strconv.Atoi(c.Query("shop_id"))

	if err != nil {
		response.ValidateFail(c, "参数错误")
	}

	resp := response.CouponResp{}

	resp.CouponList, err = service.CouponService.List(0, int32(shopId))
	resp.UserCouponList, err = service.UserCouponService.List(int32(userId), 0, int32(shopId), model.UserCouponStatusNormal)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, resp)
}
