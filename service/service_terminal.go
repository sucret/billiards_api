package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"billiards/request"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type terminalService struct {
	db *gorm.DB
}

var TerminalService = &terminalService{
	db: mysql.GetDB(),
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
	if err = e.db.Where("terminal_id = ?", form.TerminalId).Find(&terminal).Error; err != nil {
		err = errors.New("设备不存在")
		return
	}

	if terminal.Type == model.TerminalTypePicReader {
		err = errors.New("该设备不支持此类操作")
		return
	}

	terminal.Status = form.Status
	err = e.db.Save(&terminal).Error

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
