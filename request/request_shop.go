package request

import "billiards/pkg/mysql/model"

type ShopList struct {
	Page int32 `form:"page" json:"page"`
}

type SaveShop struct {
	model.Shop
}

// 自定义错误信息
func (saveShop SaveShop) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":    "店铺名称不能为空",
		"address.required": "地址不能为空",
		//"shopkeeper.required": "联系人不能为空",
		//"mobile.required":     "手机号码不能为空",
		//"mobile.mobile":       "手机号码格式不正确",
		//"region_id.required":  "地区为必选项",
	}
}
