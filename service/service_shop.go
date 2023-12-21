package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/pkg/ws_entity"
	"billiards/request"
	"billiards/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
)

type shopService struct {
	db    *gorm.DB
	redis *redis.Client
}

var ShopService = &shopService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

var socketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 店铺状态的通道
//var shopStatusChan = make(map[int32]chan int32)

// 推送店铺变更信号
// 如果没有shopStatusChan[shopId] 则表示没有连接socket，直接返回
func (s *shopService) PushShopStatusChan(shopId int32) {
	ws_entity.PushShopStatusChan(shopId)

}

// 店铺状态socket
func (s *shopService) StatusSocket(c *gin.Context) {
	sId, err := strconv.Atoi(c.Query("shop_id"))
	shopId := int32(sId)

	// 将当前http连接升级为websocket连接
	conn, err := socketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
		fmt.Println("socket close")
		// 释放shop chan
		//delete(shopStatusChan, shopId)
		ws_entity.DeleteShopStatusChan(shopId)
	}(conn)

	// 获取shop chan
	ws_entity.InitShopStatusChan(shopId)

	sendStatusMsg := func(shopId int32) {
		shopInfo, err := s.ShopStatus(shopId)
		resp, _ := json.Marshal(shopInfo)
		err = conn.WriteMessage(websocket.TextMessage, resp)
		if err != nil {
			return
		}
	}

	sendStatusMsg(shopId)

	// 异步监听通道消息，有消息就推送店铺状态给客户端
	go func() {
		ws_entity.ConsumeShopStatusChan(shopId, func(shopId int32) {
			sendStatusMsg(shopId)
		})
		//for {
		//	shId := <-shopChan
		//	fmt.Println("监听到店铺状态变更，推送消息...")
		//	sendStatusMsg(shId)
		//}
	}()

	// 建立一个映射店铺id的chan类型的map
	// 如果店铺有变动就往chan中推送店铺的信息
	// 在这里监控chan
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("收到消息:%s \n", msg)

		err = conn.WriteMessage(websocket.TextMessage, []byte("已收到消息"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 获取店铺状态信息
func (s *shopService) ShopStatus(shopId int32) (resp response.ShopStatusResp, err error) {
	shop := model.Shop{}
	err = s.db.Preload("TableList").
		Where("shop_id = ?", shopId).
		First(&shop).Error

	for _, v := range shop.TableList {
		resp.TableStatusList = append(resp.TableStatusList, response.TableStatus{TableID: v.TableID, Status: v.Status})
	}

	return
}

func (s *shopService) Detail(shopId int) (shop model.Shop, err error) {
	err = s.db.Preload("TableList").
		Preload("TableList.TerminalList").
		Where("shop_id = ?", shopId).
		First(&shop).Error

	if err != nil {
		err = errors.New("店铺不存在")
	}
	return
}

func (s *shopService) List() (list []*response.Shop) {
	s.db.Order("shop_id DESC").
		Preload("TableList").
		Find(&list)

	for _, v := range list {
		v.BilliardsTableNum = len(v.TableList)

		for _, t := range v.TableList {
			if t.Status == model.TableStatusClose {
				v.BilliardsTableFreeNum = v.BilliardsTableFreeNum + 1
			}
		}
		v.BilliardsPrice = 0
		if len(v.TableList) > 0 {
			v.BilliardsPrice = v.TableList[0].Price
		}
	}

	return
}

func (s *shopService) ListWithDistance(lat, lng float64) (shopList []*response.Shop) {
	shopList = s.List()

	if lat > 0 && lng > 0 {
		for _, v := range shopList {
			if v.Latitude > 0 && v.Longitude > 0 {
				v.Distance = math.Round(tool.Distance(lat, lng, v.Latitude, v.Longitude)/100) / 10
			}
		}
	}

	return
}

func (s *shopService) Save(form request.SaveShop) (shop model.Shop, err error) {

	if form.ShopID != 0 {
		if err = s.db.Where("shop_id = ?", form.ShopID).First(&shop).Error; err != nil {
			err = errors.New("店铺不存在")
			return
		}

		shop.Name = form.Name
		shop.Status = form.Status
		shop.Address = form.Address
		if err = s.db.Save(&shop).Error; err != nil {
			return
		}
	} else {
		shop = model.Shop{
			Name:    form.Name,
			Status:  form.Status,
			Address: form.Address,
		}
		err = s.db.Create(&shop).Error
	}
	return
}
