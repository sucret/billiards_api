package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type tableOrderApi struct{}

var TableOrderApi = new(tableOrderApi)

func (*tableOrderApi) Terminate(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	tableOrderId, err := strconv.Atoi(c.Query("table_order_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	order, err := service.TableOrderService.Terminate(userId, tableOrderId)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (*tableOrderApi) Detail(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	tableOrderId, err := strconv.Atoi(c.Query("table_order_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	order, err := service.TableOrderService.Detail(userId, tableOrderId)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (*tableOrderApi) List(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))
	orderType, err := strconv.Atoi(c.Query("type"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	list, err := service.TableOrderService.List(userId, orderType)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, list)
}

func (*tableOrderApi) Create(c *gin.Context) {
	//var form request.OrderCreate
	//if err := c.ShouldBindJSON(&form); err != nil {
	//	response.ValidateFail(c, request.GetErrorMsg(form, err))
	//	return
	//}
	//userId, _ := strconv.Atoi(c.GetString("userId"))
	//
	//order, err := service.TableOrderService.Create(form.TableID, int32(userId))
	//if err != nil {
	//	response.BusinessFail(c, err.Error())
	//	return
	//}
	//
	//response.Success(c, order)
}

// 查询用户当前订单的支付结果
func (*tableOrderApi) PayResult(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))

	orderId, _ := strconv.Atoi(c.Query("order_id"))
	if orderId == 0 {
		response.BusinessFail(c, "参数错误")
		return
	}

	status, err := service.TableOrderService.GetPayStatus(orderId, int32(userId))
	if err != nil {
		response.BusinessFail(c, "订单未支付")
		return
	}

	response.Success(c, status)
}
