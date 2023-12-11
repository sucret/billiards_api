package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"billiards/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type userCouponService struct {
	db *gorm.DB
}

var UserCouponService = &userCouponService{
	db: mysql.GetDB(),
}

func (u *userCouponService) disable(db *gorm.DB, userCouponId int32) (err error) {
	userCoupon := model.UserCoupon{}

	if err = db.Where("user_coupon_id = ? AND status = ?", userCouponId, model.UserCouponStatusNormal).
		First(&userCoupon).Error; err != nil {
		err = errors.New("未查询到优惠券信息")
		return
	}

	userCoupon.Status = model.UserCouponStatusCancel
	err = db.Save(&userCoupon).Error
	return
}

func (u *userCouponService) GetByID(userCouponId int32) (userCoupon model.UserCoupon, err error) {
	if err = u.db.Where("user_coupon_id = ?", userCouponId).
		First(&userCoupon).Error; err != nil {
		err = errors.New("未查到该优惠券")
		return
	}

	return
}

// 使用优惠券
func (u *userCouponService) Use(db *gorm.DB, userCouponId, userId int32) (userCoupon model.UserCoupon, err error) {
	if err = db.Where("user_coupon_id = ? AND status = ? AND user_id = ?", userCouponId, model.UserCouponStatusNormal, userId).
		First(&userCoupon).Error; err != nil {
		err = errors.New("未查到该优惠券")
		return
	}

	userCoupon.Status = model.UserCouponStatusUsed
	if err = db.Save(&userCoupon).Error; err != nil {
		return
	}

	return
}

func (u *userCouponService) Buy(userId, couponId int32) (resp *response.CouponOrderPrePayParam, err error) {
	coupon, err := CouponService.GetById(couponId)
	if err != nil {
		return
	}

	// 根据价格生成付款单
	order, err := PaymentService.MakeCouponOrder(u.db, coupon.Price, model.POTypeCoupon, 1, userId, coupon.Name)
	if err != nil {
		return nil, err
	}
	fmt.Println(order)
	return
}

func (u *userCouponService) PaySuccess() {

}

// 优惠券列表
func (u *userCouponService) List(userId, couponId int32, status int) (list []model.UserCoupon, err error) {
	query := u.db.Preload("Coupon")

	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if couponId > 0 {
		query = query.Where("coupon_id = ?", couponId)
	}

	if err = query.Order("user_coupon_id DESC").Find(&list).Error; err != nil {
		return
	}

	return
}

// 发放优惠券
func (u *userCouponService) HandOut(userId, couponId int32) (userCoupon model.UserCoupon, err error) {
	// 获取优惠券信息
	coupon, err := CouponService.GetById(couponId)
	if err != nil {
		return
	}

	userCoupon = model.UserCoupon{
		CouponID: coupon.CouponID,
		UserID:   userId,
		Status:   model.UserCouponStatusNormal,
		ShopID:   coupon.ShopID,
	}

	u.db.Save(&userCoupon)

	return
}

func (u *userCouponService) UseByCouponID(db *gorm.DB, couponId, userId int32) (userCoupon model.UserCoupon, err error) {
	if err = db.Where("coupon_id = ? AND status = ? AND user_id = ?", couponId, model.UserCouponStatusNormal, userId).
		First(&userCoupon).Error; err != nil {
		err = errors.New("当前用户没有未使用的优惠券")
		return
	}

	userCoupon.Status = model.UserCouponStatusUsed
	if err = db.Save(&userCoupon).Error; err != nil {
		return
	}

	return
}
