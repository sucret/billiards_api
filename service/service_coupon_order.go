package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type couponOrderService struct {
	db *gorm.DB
}

var CouponOrderService = &couponOrderService{
	db: mysql.GetDB(),
}

func (*couponOrderService) PaySuccess(db *gorm.DB, order *model.CouponOrder) (err error) {
	// 更新优惠券状态
	order.Status = 2
	if err = db.Save(&order).Error; err != nil {
		err = errors.New("更新优惠券订单失败")
		return
	}

	// 发放优惠券
	_, err = UserCouponService.HandOut(order.UserID, order.CouponID)
	if err != nil {
		return err
	}
	return
}

// 创建优惠券子订单
func (*couponOrderService) Create(tx *gorm.DB, couponId, userId, orderId int32) (
	couponOrder model.CouponOrder, err error) {

	// 1、判断优惠券是否有效
	coupon, err := CouponService.GetById(couponId)
	if err != nil {
		return
	}

	// 2、创建子订单
	couponOrder = model.CouponOrder{
		OrderID:   orderId,
		UserID:    userId,
		CouponID:  couponId,
		Status:    model.CouponOrderStatusDefault,
		PayAmount: coupon.Price,
	}

	if err = tx.Save(&couponOrder).Error; err != nil {
		log.GetLogger().Error("make_coupon_order_error", zap.String("msg", err.Error()))
	}

	return
}
