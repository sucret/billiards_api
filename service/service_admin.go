package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"
	redis_ "gin-api/pkg/redis"
	"gin-api/request"
	"gin-api/response"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type adminService struct {
	db    *gorm.DB
	redis *redis.Client
}

var AdminService = &adminService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

func (a *adminService) Save(data request.SaveAdmin) (admin model.Admin, err error) {
	admin.AdminID = data.AdminId
	admin.Username = data.Username
	admin.Mobile = data.Mobile
	admin.Nickname = data.Nickname

	if data.Password != "" {
		h := md5.New()
		h.Write([]byte(data.Password))

		admin.Password = hex.EncodeToString(h.Sum(nil))
	}

	tx := a.db.Begin()

	if data.AdminId > 0 {
		if err = a.db.Where("admin_id = ?", data.AdminId).First(&model.Admin{}).Error; err != nil {
			err = errors.New("用户不存在")
			return
		}

		a.db.Model(&model.Admin{}).Where("admin_id = ?", data.AdminId).Updates(&admin)

		// 删除用户对应的角色
		err = a.db.Where("admin_id = ?", data.AdminId).Delete(model.AdminRole{}).Error
		if err != nil {
			tx.Rollback()
			return
		}
	} else {
		err = a.db.Create(&admin).Error
	}

	// 设置角色信息
	for _, val := range data.Role {
		adminRole := model.AdminRole{
			AdminID: admin.AdminID,
			RoleID:  val,
		}

		err = a.db.Create(&adminRole).Error

		if err != nil {
			err = errors.New("保存失败")
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

func (a *adminService) Detail(adminId uint) (admin response.AdminDetail, err error) {
	a.db.Where("admin_id = ?", adminId).Find(&admin.Admin)

	a.db.Table("admin_role").
		Select("admin_role.role_id, role.name").
		Joins("left join role on admin_role.role_id = role.role_id").
		Where("admin_role.admin_id = ?", adminId).
		Scan(&admin.RoleInfo)

	// 获取所有角色
	a.db.Select("role_id, name").Find(&admin.RoleList)

	admin.NodeList = NodeService.Tree([]int{1, 2, 3})

	return
}

func (a *adminService) List() (list []model.Admin, err error) {
	err = a.db.Preload("RoleList").Order("admin_id DESC").Find(&list).Error
	return
}

// AdminLogin 登陆
func (a *adminService) AdminLogin(params request.AdminLogin) (admin *model.Admin, err error) {
	if err = a.db.Where("mobile = ?", params.Mobile).First(&admin).Error; err != nil {
		err = errors.New("用户不存在")
	}

	fmt.Println("form", params)
	if params.Captcha != "" {
		// 判断验证码是否正确
		if code, _ := a.redis.Get("sms-code_" + params.Mobile).Result(); code != params.Captcha {
			err = errors.New("验证码不正确")
			return
		}
	} else if params.Password != "" {
		h := md5.New()
		h.Write([]byte(params.Password))

		if hex.EncodeToString(h.Sum(nil)) != admin.Password {
			err = errors.New("密码不正确")
			return
		}
	} else {
		err = errors.New("验证码或密码不正确")
	}

	return
}

func (a *adminService) GetAdminById(adminId int) (admin *model.Admin, err error) {
	if err = a.db.Where("admin_id = ?", adminId).First(&admin).Error; err != nil {
		err = errors.New("用户不存在")
	}

	return
}

// GetLoginSms 获取登陆验证码
func (a *adminService) GetLoginSms(mobile string) (code string, err error) {
	var admin model.Admin

	if err = a.db.Where("mobile = ?", mobile).First(&admin).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < 4; i++ {
		code = fmt.Sprintf("%s%d", code, rand.Intn(10))
	}

	// TODO 发送验证码

	// 验证码写入redis 120秒后过期
	if err = a.redis.Set("sms-code_"+mobile, code, 120*time.Second).Err(); err != nil {
		err = errors.New("验证码写入失败")
	}

	return
}

// CheckAdminPermission 验证后端管理员权限
func (a *adminService) CheckAdminPermission(userId int, path string) bool {
	// 超级管理员可以访问所有接口
	if userId == 1 {
		return true
	}

	admin := model.Admin{}
	a.db.Preload("RoleList").
		Preload("RoleList.NodeList").
		Where("admin_id = ?", userId).First(&admin)

	for _, role := range admin.RoleList {
		for _, node := range role.NodeList {
			if node.Path == path {
				return true
			}
		}
	}

	return false
}
