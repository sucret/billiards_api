package request

type SaveNode struct {
	NodeId       int32  `form:"node_id" json:"node_id"`
	Title        string `form:"title" json:"title" binding:"required"`
	Name         string `form:"name" json:"name"`
	ParentNodeId int32  `form:"parent_node_id" json:"parent_node_id"`
	Path         string `form:"path" json:"path"`
	Type         int    `form:"type" json:"type"`
	Icon         string `form:"icon" json:"icon"`
	Redirect     string `form:"redirect" json:"redirect"`
	Component    string `form:"component" json:"component"`
}

func (SaveNode) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"title.required": "节点名称为必选项",
		// "component.required": "组件名称为必选项",
	}
}
