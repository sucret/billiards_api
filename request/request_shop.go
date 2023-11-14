package request

type ShopList struct {
	Page int32 `form:"page" json:"page"`
}

type SaveShop struct {
	ShopId     uint   `form:"shop_id" json:"shop_id"`
	Name       string `form:"name" json:"name" binding:"required"`
	Address    string `form:"address" json:"address" binding:"required"`
	RegionId   uint   `form:"region_id" json:"region_id" binding:"required"`
	Shopkeeper string `form:"shopkeeper" json:"shopkeeper" binding:"required"`
	Mobile     string `form:"mobile" json:"mobile" binding:"required,mobile"`
}

// 自定义错误信息
func (saveShop SaveShop) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":       "用户名称不能为空",
		"address.required":    "地址不能为空",
		"shopkeeper.required": "联系人不能为空",
		"mobile.required":     "手机号码不能为空",
		"mobile.mobile":       "手机号码格式不正确",
		"region_id.required":  "地区为必选项",
	}
}
