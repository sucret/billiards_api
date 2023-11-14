package request

type SaveStation struct {
	ShopId    uint   `form:"shop_id" json:"shop_id" binding:"required"`
	StationId uint   `form:"station_id" json:"station_id"`
	Camera1   string `form:"camera1" json:"camera1" binding:"required"`
	Camera2   string `form:"camera2" json:"camera2" binding:"required"`
}

func (saveStation SaveStation) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"shop_id.required": "店铺为必选项",
		"camera1.required": "请输入摄像头参数",
		"camera2.required": "请输入摄像头参数",
	}
}
