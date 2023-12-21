package model

import (
	"billiards/pkg/ws_entity"
	"gorm.io/gorm"
	"time"
)

const (
	TableStatusOpen  = iota + 1 // 开启
	TableStatusClose            // 关闭
	TableStatusStop
)

func (t *Table) AfterUpdate(tx *gorm.DB) (err error) {
	// 这里需要延迟去推送消息，不然事务没提交消息就被消费掉了
	go func() {
		time.Sleep(1 * time.Second)
		ws_entity.PushShopStatusChan(t.ShopID)
	}()

	return
}
