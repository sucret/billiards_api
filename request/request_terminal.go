package request

type ChangeTerminalStatus struct {
	TerminalId int32 `form:"terminal_id" json:"terminal_id" binding:"required"`
	Status     int32 `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (ChangeTerminalStatus) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"terminal_id.required": "设备ID为必传",
		"status.required":      "设备状态为必传",
		"status.oneof":         "状态值只能是1或者2",
	}
}

type SaveTerminal struct {
	ShopID     int32 `form:"shop_id" json:"shop_id" binding:"required"`
	TableID    int32 `form:"table_id" json:"table_id" binding:"required"`
	TerminalId int32 `form:"terminal_id" json:"terminal_id"`
	//Status      int32  `form:"status" json:"status" binding:"required,oneof=1 2"`
	Type int32  `form:"type" json:"type" binding:"required,oneof=1 2 3 5"`
	URL  string `form:"url" json:"url" binding:"required"`
}

func (SaveTerminal) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"shop_id.required":  "店铺ID为必传",
		"table_id.required": "球桌ID为必传",
		"status.required":   "设备状态为必传",
		"status.oneof":      "状态值只能是1、2",
		"type.required":     "设备类型为必选项",
		"type.oneof":        "设备类型值只能是1、2、3、5",
		"url.required":      "操作地址为必填项",
	}
}
