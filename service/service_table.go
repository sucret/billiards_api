package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/request"
	"errors"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"sync"
)

type tableService struct {
	db    *gorm.DB
	redis *redis.Client
	lock  sync.Mutex
}

var TableService = &tableService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
	lock:  sync.Mutex{},
}

func (t *tableService) Save(form request.SaveTable) (table model.Table, err error) {
	if form.TableID != 0 {
		if err = t.db.Where("table_id = ?", form.TableID).First(&table).Error; err != nil {
			err = errors.New("球桌不存在")
			return
		}

		table.Name = form.Name
		if err = t.db.Save(&table).Error; err != nil {
			return
		}
	} else {
		// todo 新增球桌之后需要生成二维码

		table = model.Table{
			Name:   form.Name,
			ShopID: form.ShopID,
			Status: model.TableStatusClose,
		}
		err = t.db.Create(&table).Error
	}

	return
}

func (t *tableService) Activate(tableId int) (table model.Table, err error) {
	// 操作设备加锁
	t.lock.Lock()
	defer t.lock.Unlock()

	tx := t.db.Begin()

	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Where("table_id = ?", tableId).
		First(&table).Error; err != nil {

		err = errors.New("设备不存在")
		tx.Rollback()
		return
	}

	// 调用终端服务去开启，成功之后修改表

	table.Status = model.TableStatusOpen
	if err = tx.Save(table).Error; err != nil {
		err = errors.New("开启失败")
		tx.Rollback()
		return
	}

	tool.Dump(table)

	tx.Commit()
	return
}
