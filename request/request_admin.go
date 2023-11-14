package request

type SaveAdmin struct {
	AdminId  int32   `json:"admin_id"`
	Username string  `json:"username" binding:"required"`
	Nickname string  `json:"nickname"`
	Mobile   string  `json:"mobile" binding:"required,mobile"`
	Password string  `json:"password"`
	Role     []int32 `json:"role"`
}

func (SaveAdmin) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"username.required": "用户名不能为空",
		"mobile.required":   "手机号码不能为空",
		"mobile.mobile":     "手机号码格式不正确",
	}
}

type AdminLogin struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password"`
	Captcha  string `form:"captcha" json:"captcha"`
}

func (adminLogin AdminLogin) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"mobile.required": "手机号码不能为空",
		"mobile.mobile":   "手机号码格式不正确",
	}
}

type AdminSendLoginSms struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
}

func (adminSendLoginSms AdminSendLoginSms) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"mobile.required": "手机号码不能为空",
		"mobile.mobile":   "手机号码格式不正确",
	}
}
