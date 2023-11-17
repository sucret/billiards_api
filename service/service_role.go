package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/request"
	"errors"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type roleService struct {
	db    *gorm.DB
	redis *redis.Client
}

var RoleService = &roleService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

// 列表
func (r *roleService) List(page int) (roleList []model.Role) {
	r.db.Order("role_id desc").Find(&roleList)
	return
}

func (r *roleService) Save(param request.SaveRole) (role model.Role, err error) {
	if param.RoleID == 1 {
		err = errors.New("超级管理员不可编辑")
		return
	}

	tx := r.db.Begin()

	if param.RoleID > 0 {
		if err = tx.Where("role_id = ?", param.RoleID).First(&role).Error; err != nil {
			err = errors.New("角色不存在")
			tx.Rollback()
			return
		}

		role.Name = param.Name
		role.Status = param.Status
		if err = tx.Save(&role).Error; err != nil {
			tx.Rollback()
			return
		}

	} else {
		role = model.Role{
			Name:   param.Name,
			Status: param.Status,
		}
		err = tx.Create(&role).Error
	}

	// 写入权限信息
	err = r.SaveNode(role.RoleID, param.NodeIdList)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (r *roleService) Detail(roleID uint) (role model.Role, err error) {
	if err = r.db.Preload("NodeList").Where("role_id = ?", roleID).First(&role).Error; err != nil {
		err = errors.New("角色不存在")
	}

	return
}

func (r *roleService) SaveNode(roleID int32, nodeIds []int32) (err error) {
	// 删除当前角色所有的权限

	tx := r.db.Begin()

	err = r.db.Where("role_id = ?", roleID).Delete(model.RoleNode{}).Error
	if err != nil {
		tx.Rollback()
		return
	}

	for _, val := range nodeIds {
		roleNode := model.RoleNode{
			RoleID: roleID,
			NodeID: val,
		}

		err = r.db.Create(&roleNode).Error

		if err != nil {
			err = errors.New("权限写入错误")
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

func (r *roleService) GetRoleNode(roleId int32) (nodeList []int32, err error) {
	var node []model.RoleNode
	r.db.Where("role_id = ?", roleId).Find(&node)

	for _, val := range node {
		nodeList = append(nodeList, val.NodeID)
	}
	return
}
