package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"sync"
	"time"
)

type orderService struct {
	db    *gorm.DB
	redis *redis.Client
	lock  sync.Mutex
}

var OrderService = &orderService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
	lock:  sync.Mutex{},
}

// 终止订单
//
// 检查球杆是否归还
// 检查球是否够
// 关闭设备
// 结算
// 发起微信退款
// 修改订单状态
func (o *orderService) Terminate(userId, orderId int) (order model.Order, err error) {
	tx := o.db.Begin()

	//if err = tx.Set("gorm:query_option", "FOR UPDATE").First(&order, orderId).Error; err != nil {
	//	tx.Rollback()
	//	return
	//}

	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		//Preload("Table").
		Where("user_id = ? AND order_id = ? AND status = ?", userId, orderId, model.OrderStatusPaySuccess).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
		tx.Rollback()
		return model.Order{}, err
	}

	// 关闭球桌
	_, err = TableService.Disable(order.TableID)
	if err != nil {
		tx.Rollback()
		return
		//return model.Order{}, err
	}

	// 修改订单状态

	order.Status = model.OrderStatusFinised

	if err = tx.Save(&order).Error; err != nil {
		tx.Rollback()
		err = errors.New("关闭失败")
		return
	}

	tx.Commit()
	tool.Dump(order)

	return
}

func (o *orderService) Detail(userId, orderId int) (order model.Order, err error) {
	if err := o.db.Preload("Table").
		Preload("Table.Shop").
		Where("user_id = ? AND order_id = ?", userId, orderId).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
	}

	return
}

func (o *orderService) List(userId int, orderType int) (list []model.Order, err error) {

	if err := o.db.
		Preload("Table").
		Preload("Table.Shop").
		Where("status = ? AND user_id = ?", orderType, userId).
		Order("order_id desc").
		Find(&list).Error; err != nil {
		fmt.Println(err)
	}

	return
}

// 查询用户自己订单的支付结果
func (o *orderService) GetPayStatus(orderNum string, userId int32) (successful bool, err error) {
	order := model.Order{}
	if err = o.db.Where("order_num = ? AND user_id = ?", orderNum, userId).
		First(&order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	successful = order.Status == model.OrderStatusPaySuccess

	return
}

// 取消10分钟前未支付的订单
func (o *orderService) TimingCancel() {
	var orderList []model.Order

	t, _ := time.ParseDuration("-10m")
	o.db.Where("status = ? AND created_at < ?", model.OrderStatusDefault,
		time.Now().Add(t).Format("2006-01-02 15:04:05")).
		Find(&orderList)

	for _, v := range orderList {
		fmt.Println(v)
		v.Status = model.OrderStatusAutoCancel
		err := o.db.Save(&v).Error
		if err != nil {
			log.GetLogger().Error("auto cancel error",
				zap.Any("order", v),
				zap.String("error", err.Error()))
		}
	}
}

// 生成订单号
func (o *orderService) GenerateOrderNum() (orderNum string) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(999999)
	orderNum = time.Now().Format("20060102150405") + fmt.Sprintf("%0*d", 6, num)

	return
}

// 创建订单
func (o *orderService) Create(tableId, userId int32) (order model.Order, err error) {
	o.GenerateOrderNum()
	// todo 暂时先在创建订单的时候关闭过期的订单，后边需要移动到定时任务里边去
	go o.TimingCancel()

	// todo 这里的锁需要再斟酌一下，如果开启的话，那同时所有的球桌就只能有一个人下单（单服务器的情况下）
	o.lock.Lock()
	defer o.lock.Unlock()

	tx := o.db.Begin()

	// 先查询出球桌，判断一下是否可以开台
	table := model.Table{}
	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Where("table_id = ? AND status = ?", tableId, model.TableStatusClose).
		First(&table).Error; err != nil {

		err = errors.New("当前球桌不可用，请更换其它球桌")
		tx.Rollback()
		return
	}

	// 需要判断这个球桌是否有支付中的订单(当前扫码用户的订单排除掉）
	tmpOrder := model.Order{}
	tx.Where("table_id = ? AND status = ? AND user_id !=", tableId, model.OrderStatusDefault, userId).
		First(&tmpOrder)
	if tmpOrder.OrderID > 0 {
		err = errors.New("当前球桌不可用，请更换其它球桌")
		tx.Rollback()
		return
	}

	// 创建订单
	order = model.Order{
		ShopID:   table.ShopID,
		TableID:  table.TableID,
		Status:   model.OrderStatusDefault,
		UserID:   userId,
		OrderNum: o.GenerateOrderNum(),
	}

	if err = tx.Create(&order).Error; err != nil {
		log.GetLogger().Error("create_order",
			zap.String("err_msg", err.Error()),
			zap.Any("order", order))

		err = errors.New("创建订单，请重试")
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

func (o *orderService) PaySuccess(orderNum string) (order model.Order, err error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	tx := o.db.Begin()
	err = tx.Where("order_num = ? AND status = ?", orderNum, model.OrderStatusDefault).
		First(&order).Error
	if err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("err_msg", "没有查到订单或者订单支付状态已经完成"),
			zap.String("order_num", orderNum))
		tx.Rollback()
		return
	}

	tool.Dump(order)

	order.Status = model.OrderStatusPaySuccess
	order.PaidAt = model.Time(time.Now())

	if err = tx.Save(&order).Error; err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("msg", "更新订单状态失败"),
			zap.String("err_msg", err.Error()),
			zap.String("order_num", orderNum))
		tx.Rollback()
	}

	// 开台
	table, err := TableService.Enable(order.TableID)
	if err != nil {
		tx.Rollback()

		log.GetLogger().Error("pay_notify",
			zap.String("msg", "开台失败"),
			zap.Any("table", table),
			zap.String("order_num", orderNum))
		tx.Rollback()

		return model.Order{}, err
	}

	tx.Commit()
	return
}
