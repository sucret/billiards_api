package ws_entity

import "fmt"

// 店铺状态的通道
var shopStatusChan = make(map[int32]chan int32)

func InitShopStatusChan(shopId int32) {
	_, ok := shopStatusChan[shopId]
	if !ok {
		shopChan := make(chan int32, 10)
		shopStatusChan[shopId] = shopChan
	}
}

// 推送店铺变更信号
// 如果没有shopStatusChan[shopId] 则表示没有连接socket，直接返回
func PushShopStatusChan(shopId int32) {
	_, ok := shopStatusChan[shopId]
	if !ok {
		return
	}

	shopStatusChan[shopId] <- shopId
}

func ConsumeShopStatusChan(shopId int32, fn func(int32)) {
	shopChan, ok := shopStatusChan[shopId]
	if !ok {
		return
	}

	for {
		sId := <-shopChan
		fmt.Println("监听到店铺状态变更，推送消息...")
		fn(sId)
	}
}

func DeleteShopStatusChan(shopId int32) {
	delete(shopStatusChan, shopId)
}
