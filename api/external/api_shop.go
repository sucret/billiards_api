package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type shopApi struct{}

var ShopApi = new(shopApi)

// 获取店铺设备状态
func (*shopApi) TerminalStatus(c *gin.Context) {
	shopId, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	resp, err := service.ShopService.ShopStatus(int32(shopId))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, resp)
}

func (*shopApi) TerminalStatusSocket(c *gin.Context) {
	service.ShopService.StatusSocket(c)
}

func (*shopApi) ChanTest(c *gin.Context) {
	shopId, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}
	service.ShopService.PushShopStatusChan(int32(shopId))
}
