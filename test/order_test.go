package test

import (
	"billiards/pkg/mysql"
	"billiards/pkg/tool"
	"billiards/request"
	"billiards/service"
	"fmt"
	"testing"
)

func TestOrderDetail(t *testing.T) {
	_, _ = service.OrderService.Detail(60, 3)
}

// 测试创建优惠券
func TestCreateOrder(t *testing.T) {
	var form request.OrderCreate
	// 开台订单
	form = request.OrderCreate{
		TableID: 1,
		//CouponID: 3,
	}

	// 充值订单
	//form = request.OrderCreate{
	//	IsRecharge:         true,
	//	RechargeAmountType: 1,
	//}

	// 购买优惠券
	//form = request.OrderCreate{
	//	CouponID: 6,
	//}

	order, err := service.OrderService.Create(3, form)
	if err != nil {
		fmt.Println(err)
		return
	}

	tool.Dump(order)
}

func TestOrderSuccess(t *testing.T) {
	db := mysql.GetDB()
	err := service.OrderService.PaySuccess(db, 47)
	if err != nil {
		fmt.Println(err)
		return
	}
}
