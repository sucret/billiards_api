package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"billiards/request"
	"errors"
	"gorm.io/gorm"
)

type couponService struct {
	db *gorm.DB
}

var CouponService = &couponService{
	db: mysql.GetDB(),
}

func (c *couponService) List(couponId int32, shopId int32) (list []model.Coupon, err error) {
	query := c.db
	if couponId > 0 {
		query = query.Where("coupon_id = ?", couponId)
	}

	if shopId > 0 {
		query = query.Where("shop_id = ?", shopId)
	}

	err = query.Order("coupon_id DESC").Find(&list).Error

	return
}

func (c *couponService) GetById(couponId int32) (coupon *model.Coupon, err error) {
	if err = c.db.Where("coupon_id = ?", couponId).First(&coupon).Error; err != nil {
		err = errors.New("优惠券不存在")
	}
	return
}

func (c *couponService) Save(form request.SaveCoupon) (coupon model.Coupon, err error) {
	if form.CouponID > 0 {
		if err = c.db.Where("coupon_id = ?", form.CouponID).First(&coupon).Error; err != nil {
			err = errors.New("优惠券没找到")
			return
		}
	}

	coupon.Price = form.Price
	coupon.Type = form.Type
	coupon.Name = form.Name
	coupon.ShopID = form.ShopID
	coupon.Duration = form.Duration
	coupon.TableType = form.TableType

	if err = c.db.Save(&coupon).Error; err != nil {
		err = errors.New("更新优惠券失败")
		return
	}

	return
}
