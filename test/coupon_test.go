package test

import (
	"billiards/pkg/mysql/model"
	"billiards/pkg/tool"
	"billiards/request"
	"billiards/service"
	"math/rand"
	"testing"
	"time"
)

// 测试创建优惠券
func TestCreateCoupon(t *testing.T) {
	form := request.SaveCoupon{
		Name:      "门店优惠券",
		Price:     100,
		ShopID:    1,
		Duration:  60,
		Type:      model.CouponTypeDuration,
		TableType: model.CouponTableTypeGeneral,
	}

	t.Run("create", func(t *testing.T) {
		save, err := service.CouponService.Save(form)
		if err != nil {
			return
		}

		tool.Dump(save)
	})

	t.Run("update", func(t *testing.T) {
		form.CouponID = 3
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(100)

		form.Duration = int32(num)
		save, err := service.CouponService.Save(form)
		if err != nil {
			return
		}

		tool.Dump(save)
	})
}
