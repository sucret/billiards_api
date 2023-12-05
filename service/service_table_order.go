package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/response"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"sync"
	"time"
)

type tableOrderService struct {
	db    *gorm.DB
	redis *redis.Client
	lock  sync.Mutex
}

var TableOrderService = &tableOrderService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
	lock:  sync.Mutex{},
}

func (o *tableOrderService) GetByTable(tableId, userId int) (order model.TableOrder, err error) {
	if err = o.db.Where("table_id = ? AND user_id = ? AND status = ?",
		tableId, userId, model.OrderStatusPaySuccess).
		First(&order).Error; err != nil {
		return
	}
	return
}

// 终止订单
//
// 检查球杆是否归还
// 检查球是否够
// 关闭设备
// 结算
// 发起微信退款
// 修改订单状态
func (o *tableOrderService) Terminate(userId, orderId int) (order *response.OrderDetail, err error) {
	// todo 判断球杆柜是否关闭

	tx := o.db.Begin()

	order = &response.OrderDetail{}
	// 获取订单信息
	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("Table").
		Preload("PaymentOrderList").
		Where("user_id = ? AND order_id = ? AND status = ?", userId, orderId, model.OrderStatusPaySuccess).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
		tx.Rollback()
		return
	}

	// 关闭球桌
	// todo
	_, err = TableService.Disable(tx, order.TableID)
	if err != nil {
		tx.Rollback()
		return
	}

	//tool.Dump(order)
	//tx.Rollback()
	//return

	// 修改订单状态
	order.Status = model.OrderStatusFinised
	order.TerminatedAt = model.Time(time.Now())

	// 结算 回写订单金额
	o.settlement(&order.TableOrder)

	// todo 退押金
	o.refund(order.TableOrder)

	if err = tx.Save(&order.TableOrder).Error; err != nil {
		tx.Rollback()
		err = errors.New("关闭失败")
		return
	}

	tx.Commit()

	// 重新获取订单信息
	order, err = o.Detail(userId, orderId)

	return
}

// 退款
func (o *tableOrderService) refund(order model.TableOrder) {
	var totalAmount int32
	for _, v := range order.PaymentOrderList {
		totalAmount = totalAmount + v.Amount
	}

	refundAmount := totalAmount - order.Amount

	PaymentService.RefundOrder(order, refundAmount)
}

func (o *tableOrderService) Detail(userId, orderId int) (order *response.OrderDetail, err error) {
	if err = o.db.Preload("Table").
		Preload("Table.Shop").
		Preload("PaymentOrderList", func(db *gorm.DB) *gorm.DB {
			return db.Select("order_id,amount").
				Where("status = ?", model.PMOStatusSuccess).
				Order("pay_mode")
		}).
		Where("user_id = ? AND order_id = ?", userId, orderId).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
	}

	o.formatClientOrder(order)

	return
}

// 计算订单的金额和时间
func (o *tableOrderService) formatClientOrder(order *response.OrderDetail) {
	var totalAmount int32
	for _, v := range order.PaymentOrderList {
		totalAmount = totalAmount + v.Amount
	}

	// 总时长 （当前支付的金额 / 单价 * 60）
	order.TotalMinutes = (totalAmount / order.Table.Price) * 60

	if order.Status == model.OrderStatusPaySuccess {
		// 使用时长
		order.UsedMinutes = int32(time.Now().Sub(time.Time(order.StartedAt)).Minutes())
	} else if order.Status == model.OrderStatusFinised {
		order.UsedMinutes = int32(time.Time(order.TerminatedAt).Sub(time.Time(order.StartedAt)).Minutes())
	}

	// 剩余时长
	order.RemainMinutes = order.TotalMinutes - order.UsedMinutes
	if order.RemainMinutes < 0 {
		order.RemainMinutes = 0
	}

	//tool.Dump(order)
}

func (o *tableOrderService) List(userId int, orderType int) (list []model.TableOrder, err error) {
	query := o.db.Preload("Table").Preload("Table.Shop")

	if orderType == 0 {
		query.Where("status in ? AND user_id = ?",
			[]int{model.OrderStatusPaySuccess, model.OrderStatusFinised}, userId)
	} else {
		query.Where("status = ? AND user_id = ?", orderType, userId)
	}
	if err := query.Order("order_id desc").
		Find(&list).Error; err != nil {
	}

	return
}

// 查询用户自己订单的支付结果
func (o *tableOrderService) GetPayStatus(orderId int, userId int32) (successful bool, err error) {
	order := model.TableOrder{}
	if err = o.db.Where("order_id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	successful = order.Status == model.OrderStatusPaySuccess

	return
}

// 取消10分钟前未支付的订单
func (o *tableOrderService) TimingCancel() {
	var orderList []model.TableOrder

	t, _ := time.ParseDuration("-10m")
	o.db.Where("status = ? AND created_at < ?", model.OrderStatusDefault,
		time.Now().Add(t).Format("2006-01-02 15:04:05")).
		Find(&orderList)

	for _, v := range orderList {
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
func (o *tableOrderService) GenerateOrderNum() (orderNum string) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(999999)
	orderNum = time.Now().Format("20060102150405") + fmt.Sprintf("%0*d", 6, num)

	return
}

// 创建订单
func (o *tableOrderService) Create(tableId, userId int32) (resp response.TableOrderPrePayParam, err error) {
	// todo 暂时先在创建订单的时候关闭过期的订单，后边需要移动到定时任务里边去
	go o.TimingCancel()

	// todo 这里的锁需要再斟酌一下，如果开启的话，那同时所有的球桌就只能有一个人下单（单服务器的情况下）
	//o.lock.Lock()
	//defer o.lock.Unlock()

	order := model.TableOrder{}
	tx := o.db.Begin()

	// 先查询出球桌，判断一下是否可以开台
	table := model.Table{}
	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("Shop").
		Where("table_id = ? AND status = ?", tableId, model.TableStatusClose).
		First(&table).Error; err != nil {

		err = errors.New("当前球桌不可用，请更换其它球桌")
		tx.Rollback()
		return
	}

	// 需要判断这个球桌是否有支付中的订单(当前扫码用户的订单排除掉）
	tmpOrder := model.TableOrder{}
	tx.Where("table_id = ? AND status = ? AND user_id !=", tableId, model.OrderStatusDefault, userId).
		First(&tmpOrder)
	if tmpOrder.OrderID > 0 {
		err = errors.New("当前球桌不可用，请更换其它球桌")
		tx.Rollback()
		return
	}

	// 创建订单
	order = model.TableOrder{
		ShopID:  table.ShopID,
		TableID: table.TableID,
		Status:  model.OrderStatusDefault,
		UserID:  userId,
	}

	if err = tx.Create(&order).Error; err != nil {
		log.GetLogger().Error("create_order", zap.String("err_msg", err.Error()),
			zap.Any("order", order))

		err = errors.New("创建订单，请重试")
		tx.Rollback()
		return
	}

	wxPayAmount := table.Shop.Deposit
	walletPayAmount := int32(0)

	// 获取用户信息
	user, _ := UserService.GetByUserId(userId)

	// 判断用户余额是否够支付，如果够，则不用微信支付
	if user.Wallet >= table.Shop.Deposit {
		wxPayAmount = 0
		walletPayAmount = table.Shop.Deposit
	} else {
		wxPayAmount = table.Shop.Deposit - user.Wallet
		walletPayAmount = user.Wallet
	}

	// 生成余额支付订单
	if walletPayAmount > 0 {
		_, err = PaymentService.MakeWalletOrder(walletPayAmount, model.POTypeTable, order.OrderID, tx)
		if err != nil {
			log.GetLogger().Error("gen_payment_order_err", zap.String("msg", err.Error()))
			tx.Rollback()
			return response.TableOrderPrePayParam{}, err
		}
	}

	// 生成预支付的参数
	if wxPayAmount > 0 {
		payment, err := PaymentService.MakeWechatPrepayOrder(
			userId, wxPayAmount, model.POTypeTable, order.OrderID, table.Name)
		if err != nil {
			log.GetLogger().Error("gen_payment_order_err", zap.String("msg", err.Error()))
			tx.Rollback()
			return response.TableOrderPrePayParam{}, err
		}
		resp.NeedWxPay = true
		resp.JsApi = payment
	}

	// 全部用钱包支付的话，订单直接改为支付成功
	if wxPayAmount == 0 {
		_, err = TableOrderService.PaySuccess(order.OrderID, tx)
	}

	resp.Order = &order

	tx.Commit()

	return
}

func (o *tableOrderService) PaySuccess(orderId int32, db *gorm.DB) (order model.TableOrder, err error) {
	err = db.Where("order_id = ? AND status = ?", orderId, model.OrderStatusDefault).
		First(&order).Error

	if err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("err_msg", "没有查到订单或者订单支付状态已经完成"),
			zap.Int32("order_id", orderId))
		return
	}

	order.Status = model.OrderStatusPaySuccess
	order.StartedAt = model.Time(time.Now())

	if err = db.Save(&order).Error; err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("msg", "更新订单状态失败"),
			zap.String("err_msg", err.Error()),
			zap.Int32("order_id", orderId))
	}

	// 开台
	table, err := TableService.Enable(order.TableID)
	if err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("msg", "开台失败"),
			zap.Any("table", table),
			zap.Int32("order_id", orderId))

		return model.TableOrder{}, err
	}

	_, err = PaymentService.MakeWalletOrderSuccess(db, orderId, order.UserID)
	if err != nil {
		log.GetLogger().Error("pay_notify", zap.String("msg", err.Error()))
		return model.TableOrder{}, err
	}

	return
}

// 结算
// 开始时间
// 结束时间
// 有没有优惠券
// 优先使用优惠券
// 剩余时间用押金结算
// 返回应该退多少押金

func (o *tableOrderService) settlement(order *model.TableOrder) {
	// 计算
	minutes := time.Time(order.TerminatedAt).Sub(time.Time(order.StartedAt)).Minutes()

	// 有优惠券的话，先减掉优惠券的时间

	order.Amount = order.Table.Price / 60 * int32(minutes)
	//order.Amount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", order.Table.Price/float64(60)*minutes), 64)

}
