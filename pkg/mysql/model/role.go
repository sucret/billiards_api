package model

import (
	"errors"
	"gorm.io/gorm"
)

func (r *Role) BeforeSave(tx *gorm.DB) (err error) {
	query := tx.Model(Role{}).Where("name = ?", r.Name)

	if r.RoleID > 0 {
		query.Where("role_id != ?", r.RoleID)
	}

	var count int64
	query.Count(&count)

	if count > 0 {
		err = errors.New("角色名称重复")
		return
	}

	return
}
