package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userService struct {
	db *gorm.DB
}

var UserService = &userService{
	db: mysql.GetDB(),
}

func (u *userService) Login(code string) (user *model.User, err error) {
	resp, err := WeApp.Login(code)
	if err != nil {
		return
	}

	if err = u.db.Where("open_id = ?", resp.OpenId).First(&user).Error; err != nil {
		// 注册新用户
		user.OpenID = resp.OpenId
		u.db.Create(user)
	}

	return
}

func (u *userService) GetByUserId(userId int32) (user model.User, err error) {
	err = u.db.Where("user_id = ?", userId).First(&user).Error
	return
}

func (u *userService) Recharge(db *gorm.DB, amount, userId int32) (user model.User, err error) {
	if err = db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		log.GetLogger().Error("recharge_error", zap.String("msg", err.Error()))
		return
	}

	user.Wallet = user.Wallet + amount

	u.db.Save(&user)
	return
}

// 用户钱包变更，如果
func (u *userService) ChangeWallet(db *gorm.DB, userId, amount int32) (err error) {
	user := model.User{}
	db.Where("user_id = ?", userId).First(&user)

	// 如果用户余额小于要扣除的金额，则支付失败
	if amount < 0 && user.Wallet < amount*-1 {
		log.GetLogger().Error("wallet_error", zap.String("msg", "余额小于支付金额，支付失败"))
		err = errors.New("余额小于支付金额，支付失败")
		return
	}

	// 扣除用户余额
	user.Wallet = user.Wallet + amount
	db.Save(&user)

	// todo 这里需要再加上用户的账户变更记录
	return
}
