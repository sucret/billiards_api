package service

import (
	"billiards/pkg/equipment"
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"billiards/request"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type terminalService struct {
	db   *gorm.DB
	lock sync.Mutex
}

var TerminalService = &terminalService{
	db:   mysql.GetDB(),
	lock: sync.Mutex{},
}

func (e *terminalService) Save(form request.SaveTerminal) (terminal model.Terminal, err error) {
	if form.TerminalId > 0 {
		if err = e.db.Where("terminal_id = ?", form.TerminalId).Find(&terminal).Error; err != nil {
			err = errors.New("设备不存在")
			return
		}
		terminal.Type = form.Type
		terminal.URL = form.URL

		if err = e.db.Save(&terminal).Error; err != nil {
			return
		}
	} else {
		terminal.ShopID = form.ShopID
		terminal.TableID = form.TableID
		terminal.Type = form.Type
		terminal.URL = form.URL
		terminal.Status = model.TerminalStatusClose
		if err = e.db.Create(&terminal).Error; err != nil {
			return
		}
	}
	return
}

func (e *terminalService) ChangeStatus(form request.ChangeTerminalStatus) (terminal model.Terminal, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	tx := e.db.Begin()

	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Where("terminal_id = ?", form.TerminalId).
		First(&terminal).Error; err != nil {
		err = errors.New("设备不存在")
		tx.Rollback()
		return
	}

	if terminal.Type == model.TerminalTypePicReader {
		err = errors.New("该设备不支持此类操作")
		tx.Rollback()
		return
	}

	if terminal.Status == form.Status {
		err = errors.New("设备状态已更新，请刷新页面")
		tx.Rollback()
		return
	}

	term := equipment.Terminal{Url: terminal.URL, Status: terminal.Status}

	log.GetLogger().Info("ChangeStatus", zap.Any("request", form))

	resp, err := term.SetStatus(form.Status)

	log.GetLogger().Info("ChangeStatus", zap.Any("response", resp))

	if err != nil {
		tx.Rollback()
		return
	}

	if resp.Status == 1 {
		terminal.Status = model.TerminalStatusOpen
	} else {
		terminal.Status = model.TerminalStatusClose
	}
	err = e.db.Save(&terminal).Error
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

func (e *terminalService) Delete(terminalId int) (err error) {
	terminal := model.Terminal{}
	if err = e.db.Where("terminal_id = ? AND deleted_at is null", terminalId).Find(&terminal).Error; err != nil {
		err = errors.New("设备不存在")
		return
	}

	time := model.Time{}

	fmt.Println(time.Local())
	terminal.DeletedAt = model.Time{}
	err = e.db.Save(&terminal).Error

	return
}
