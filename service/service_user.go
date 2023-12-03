package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
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

func (u *userService) Recharge(amount, userId int32) (user model.User, err error) {
	if err = u.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		log.GetLogger().Error("recharge_error", zap.String("msg", err.Error()))
		return
	}

	user.AccountBalance = user.AccountBalance + amount

	u.db.Save(&user)
	return
}
