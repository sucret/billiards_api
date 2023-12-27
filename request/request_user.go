package request

import "billiards/response"

// avatarUrl: "https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132"
//city: ""
//country: ""
//gender: 0
//language: ""
//nickName: "微信用户"
//province: ""

type SaveUser struct {
	Avatar   string `json:"avatar"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Nickname string `json:"nickname"`
}

type UserListReq struct {
	response.Pagination
}
