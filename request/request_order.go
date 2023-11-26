package request

type OrderCreate struct {
	TableID int32 `form:"table_id" json:"table_id" binding:"required"`
}

func (OrderCreate) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"table_id.required": "请扫描球桌上的二维码进行开台",
	}
}
