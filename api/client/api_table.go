package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type tableApi struct{}

var TableApi = new(tableApi)

func (*tableApi) Detail(c *gin.Context) {
	tableId, err := strconv.Atoi(c.Query("table_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	table, err := service.TableService.Detail(tableId)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, table)
}

// 通过桌号查询是否有订单
func (*tableApi) GetOrder(c *gin.Context) {
	tableId, err := strconv.Atoi(c.Query("table_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}
	userId, _ := strconv.Atoi(c.GetString("userId"))

	order, err := service.OrderService.GetByTable(tableId, userId)
	if err != nil {
		response.Success(c, nil)
		return
	}

	response.Success(c, order)
}
