package response

import "billiards/pkg/mysql/model"

type CouponResp struct {
	CouponList     []model.Coupon     `json:"coupon_list"`
	UserCouponList []model.UserCoupon `json:"user_coupon_list"`
}
