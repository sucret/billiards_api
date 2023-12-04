package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/request"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"sync"
	"time"
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

func (t *tableService) Detail(tableId int) (table model.Table, err error) {
	if err = t.db.Preload("Shop").
		Where("table_id = ?", tableId).
		First(&table).Error; err != nil {
		err = errors.New("查询球桌错误")
		return
	}

	return
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

func (t *tableService) Disable(tx *gorm.DB, tableId int32) (table model.Table, err error) {
	// 操作设备加锁
	//t.lock.Lock()
	//defer t.lock.Unlock()
	//
	//tx := t.db.Begin()

	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("TerminalList").
		Where("table_id = ?", tableId).
		First(&table).Error; err != nil {

		err = errors.New("球桌不存在")
		return
	}

	for _, val := range table.TerminalList {
		if val.Type == model.TerminalTypePicReader {
			continue
		}

		// todo
		//form := request.ChangeTerminalStatus{TerminalId: val.TerminalID, Status: model.TerminalStatusClose}
		//_, err = TerminalService.ChangeStatus(form)
		//if err != nil {
		//	tx.Rollback()
		//	return
		//}
	}

	table.ActivatedAt = model.Time{}
	table.Status = model.TableStatusClose
	if err = tx.Save(table).Error; err != nil {
		err = errors.New("关闭球桌失败")
		return
	}

	return
}

// 开台
func (t *tableService) Enable(tableId int32) (table model.Table, err error) {
	// 操作设备加锁
	t.lock.Lock()
	defer t.lock.Unlock()

	tx := t.db.Begin()

	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("TerminalList").
		Where("table_id = ?", tableId).
		First(&table).Error; err != nil {

		err = errors.New("球桌不存在")
		tx.Rollback()
		return
	}

	// 调用终端服务去开启，成
	//功之后修改表
	for _, val := range table.TerminalList {
		if val.Type == model.TerminalTypePicReader {
			continue
		}

		form := request.ChangeTerminalStatus{TerminalId: val.TerminalID, Status: model.TerminalStatusOpen}
		fmt.Println(form)

		// todo 开台这里先隐藏
		//_, err = TerminalService.ChangeStatus(form)
		//if err != nil {
		//	tx.Rollback()
		//	return
		//}
	}

	table.ActivatedAt = model.Time(time.Now())
	table.Status = model.TableStatusOpen
	if err = tx.Save(table).Error; err != nil {
		err = errors.New("开启失败")
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}
