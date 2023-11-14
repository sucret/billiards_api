package service

import (
	"errors"
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"
	redis_ "gin-api/pkg/redis"
	"gin-api/request"
	"gin-api/response"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type nodeService struct {
	db    *gorm.DB
	redis *redis.Client
}

var NodeService = &nodeService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

func (n *nodeService) UserNode(userId uint) (result []response.AdminNode) {
	// 获取所有菜单和页面
	var nodeList []model.Node
	n.db.Where("type IN ? AND node_id > ?", [3]int{1, 2, 3}, 3).Find(&nodeList)

	parentIdMap := make(map[int32]bool)

	for _, val := range nodeList {
		if val.ParentNodeID > 0 {
			parentIdMap[val.ParentNodeID] = true
		}
	}

	var hasChildren bool
	for _, val := range nodeList {
		if _, ok := parentIdMap[val.NodeID]; ok {
			hasChildren = true
		} else {
			hasChildren = false
		}
		if val.Type == 1 {
			hasChildren = false
		}
		if val.Icon == "" {
			val.Icon = "none"
		}
		data := response.AdminNode{
			Name:     val.Name,
			ParentId: val.ParentNodeID,
			Id:       val.NodeID,
			Meta: response.NodeMeta{
				Icon:         val.Icon,
				Show:         val.Type == 1 || val.Type == 3,
				Title:        val.Title,
				HideChildren: hasChildren && val.Type != 1,
			},
			Component: val.Component,
			Redirect:  val.Redirect,
			Path:      val.Path,
		}
		result = append(result, data)
	}
	return
}

func (n *nodeService) Delete(nodeId uint) (err error) {
	// 判断是否有子节点
	var count int64
	n.db.Model(&model.Node{}).Where("parent_node_id = ?", nodeId).Count(&count)

	if count > 0 {
		err = errors.New("有子节点，不能删除")
		return
	}

	tx := n.db.Begin()
	// 删除权限下边绑定的节点
	err = n.db.Where("node_id = ?", nodeId).Delete(model.RoleNode{}).Error
	if err != nil {
		err = errors.New("删除节点权限失败")
		tx.Rollback()
		return
	}

	err = n.db.Where("node_id = ?", nodeId).Delete(model.Node{}).Error
	if err != nil {
		err = errors.New("删除节点失败")
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

func (n *nodeService) Detail(nodeId uint) (node model.Node, err error) {
	if err = n.db.Where("node_id = ?", nodeId).First(&node).Error; err != nil {
		err = errors.New("节点不存在")
		return
	}
	return
}

func (n *nodeService) Save(param request.SaveNode) (node model.Node, err error) {
	node = model.Node{
		NodeID:       param.NodeId,
		Title:        param.Title,
		Name:         param.Name,
		ParentNodeID: param.ParentNodeId,
		Path:         param.Path,
		Type:         param.Type,
		Icon:         param.Icon,
		Redirect:     param.Redirect,
		Component:    param.Component,
	}

	// 判断名称是否重复

	if param.NodeId > 0 {
		if err = n.db.Where("node_id = ?", param.NodeId).First(&model.Node{}).Error; err != nil {
			err = errors.New("节点不存在")
			return
		}

		err = n.db.Model(&model.Node{}).Where("node_id = ?", param.NodeId).Updates(&node).Error
	} else {
		err = n.db.Create(&node).Error
	}

	return
}

func (n *nodeService) Tree(nodeType []int) (nodeTree []response.NodeTree) {
	var nodeList []model.Node

	n.db.Where("type IN ?", nodeType).Find(&nodeList)

	//n.db.Where("type IN ? AND node_id > ?", [3]int{1, 2, 3}, 3).Find(&nodeList)

	childrenList := make(map[int32][]model.Node)

	for _, value := range nodeList {
		childrenList[value.ParentNodeID] = append(childrenList[value.ParentNodeID], value)
	}

	nodeTree = getChildren(childrenList, 0)

	return
}

// 递归获取节点树
func getChildren(childrenList map[int32][]model.Node, nodeId int32) (children []response.NodeTree) {
	child, ok := childrenList[nodeId]
	if ok {
		for _, value := range child {
			data := response.NodeTree{
				NodeID:       value.NodeID,
				Title:        value.Title,
				Icon:         value.Icon,
				Type:         value.Type,
				Path:         value.Path,
				ParentNodeID: value.ParentNodeID,
				Component:    value.Component,
				Redirect:     value.Redirect,
				Children:     getChildren(childrenList, value.NodeID),
			}

			children = append(children, data)
		}
	}

	return
}
