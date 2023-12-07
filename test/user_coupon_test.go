package test

import (
	"billiards/pkg/tool"
	"billiards/service"
	"fmt"
	"testing"
)

// 发放优惠券
func TestHandOutCoupon(t *testing.T) {
	uc, err := service.UserCouponService.HandOut(3, 4)
	fmt.Println(uc, err)
}

func TestUserCouponList(t *testing.T) {
	list, err := service.UserCouponService.List(3, 3, 1)

	fmt.Println(err)
	tool.Dump(list)
}

// 直接购买优惠券
