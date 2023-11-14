package api

import (
	"gin-api/request"
	"gin-api/response"
	"gin-api/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type nodeApi struct{}

var NodeApi = new(nodeApi)

func (*nodeApi) NodeTree(c *gin.Context) {
	nodeTree := service.NodeService.Tree([]int{1, 2, 3, 4})

	response.Success(c, nodeTree)
}

// 获取所有菜单
func (*nodeApi) List(c *gin.Context) {}

func (*nodeApi) Delete(c *gin.Context) {
	nodeId, err := strconv.Atoi(c.Query("node_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	err = service.NodeService.Delete(uint(nodeId))

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, "")
}

func (*nodeApi) Save(c *gin.Context) {
	var form request.SaveNode
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	node, err := service.NodeService.Save(form)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, node)
}

func (*nodeApi) Detail(c *gin.Context) {
	nodeId, err := strconv.Atoi(c.Query("node_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	node, err := service.NodeService.Detail(uint(nodeId))

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, node)
}
