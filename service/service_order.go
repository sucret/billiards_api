package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/response"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"strconv"
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

func (o *orderService) GetByTable(tableId, userId int) (order model.Order, err error) {
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
func (o *orderService) Terminate(userId, orderId int) (order *response.OrderDetail, err error) {
	// todo 判断球杆柜是否关闭

	tx := o.db.Begin()

	order = &response.OrderDetail{}
	// 获取订单信息
	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("Table").
		Where("user_id = ? AND order_id = ? AND status = ?", userId, orderId, model.OrderStatusPaySuccess).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
		tx.Rollback()
		return
	}

	// 关闭球桌
	// todo
	_, err = TableService.Disable(order.TableID)
	if err != nil {
		tx.Rollback()
		return
	}

	// 修改订单状态
	order.Status = model.OrderStatusFinised
	order.TerminatedAt = model.Time(time.Now())

	// 结算 回写订单金额
	o.settlement(&order.Order)

	// todo 退押金

	if err = tx.Save(&order.Order).Error; err != nil {
		tx.Rollback()
		err = errors.New("关闭失败")
		return
	}

	tx.Commit()

	// 重新获取订单信息
	order, err = o.Detail(userId, orderId)

	return
}

func (o *orderService) Detail(userId, orderId int) (order *response.OrderDetail, err error) {
	if err := o.db.Preload("Table").
		Preload("Table.Shop").
		Preload("PaymentOrderList", func(db *gorm.DB) *gorm.DB {
			return db.Select("order_id,order_num,amount").Where("trade_state = ?", "SUCCESS")
		}).
		Where("user_id = ? AND order_id = ?", userId, orderId).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
	}

	o.formatClientOrder(order)

	//tool.Dump(order)

	a := time.Now().Sub(time.Time(order.PaidAt)).Minutes()
	fmt.Println(a)

	return
}

// 计算订单的金额和时间
func (o *orderService) formatClientOrder(order *response.OrderDetail) {
	var totalAmount float64
	for _, v := range order.PaymentOrderList {
		totalAmount = totalAmount + v.Amount
	}

	// 总时长 （当前支付的金额 / 单价 * 60）
	order.TotalMinutes = int32(math.Ceil((totalAmount / order.Table.Price) * 60))

	if order.Status == model.OrderStatusPaySuccess {
		// 使用时长
		order.UsedMinutes = int32(time.Now().Sub(time.Time(order.PaidAt)).Minutes())
	} else if order.Status == model.OrderStatusFinised {
		order.UsedMinutes = int32(time.Time(order.TerminatedAt).Sub(time.Time(order.PaidAt)).Minutes())
	}

	// 剩余时长
	order.RemainMinutes = order.TotalMinutes - order.UsedMinutes
	if order.RemainMinutes < 0 {
		order.RemainMinutes = 0
	}

	//tool.Dump(order)
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
func (o *orderService) Create(tableId, userId int32) (resp response.PrePayParam, err error) {
	// todo 暂时先在创建订单的时候关闭过期的订单，后边需要移动到定时任务里边去
	go o.TimingCancel()

	// todo 这里的锁需要再斟酌一下，如果开启的话，那同时所有的球桌就只能有一个人下单（单服务器的情况下）
	o.lock.Lock()
	defer o.lock.Unlock()

	order := model.Order{}
	tx := o.db.Begin()

	// 获取用户信息
	user, _ := UserService.GetByUserId(userId)

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

	// 生成预支付的参数
	payment, err := PaymentService.MakePrepayOrder(order, user.OpenID, table.Name, order.OrderNum, table.Price)

	resp.Order = &order
	resp.JsApi = payment

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

// 结算
// 开始时间
// 结束时间
// 有没有优惠券
// 优先使用优惠券
// 剩余时间用押金结算
// 返回应该退多少押金

func (o *orderService) settlement(order *model.Order) {
	// 计算
	minutes := time.Time(order.TerminatedAt).Sub(time.Time(order.PaidAt)).Minutes()

	fmt.Println("minutes", minutes)
	// 有优惠券的话，先减掉优惠券的时间

	order.Amount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", order.Table.Price/float64(60)*minutes), 64)

	tool.Dump(order)
}
