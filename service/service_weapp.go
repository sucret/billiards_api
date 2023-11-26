package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/wechat"
	"fmt"
	"gorm.io/gorm"
)

type weApp struct {
	db *gorm.DB
}

var WeApp = &weApp{
	db: mysql.GetDB(),
}

// 换取openid 判断用户是否存在
// 不存在则创建
// 生成jwt
func (w *weApp) Login(code string) (resp wechat.Code2SessionResp, err error) {
	resp, err = wechat.GetApp("client").Code2Session(code)
	if err != nil {
		return
	}

	return
}

func (w *weApp) GetAccessToken() {
	token, err := wechat.GetApp("client").GetAccessToken()
	if err != nil {
		return
	}

	fmt.Println(token)
}
