package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

var User = &user{
	db: mysql.GetDB(),
}

func (u *user) Login(code string) (user *model.User, err error) {
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
