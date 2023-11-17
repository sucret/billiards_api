package wechat

import (
	"billiards/pkg/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	ExpiresIn  int    `json:"expires_in"`
	OpenId     string `json:"openid"`
	ErrorCode  int    `json:"errcode"`
	ErrorMsg   string `json:"errmsg"`
}
type WeApp struct {
	AppId     string
	AppSecret string
}

func GetApp(t string) (instance WeApp) {
	if t == "client" {
		c := config.GetConfig().WeApp
		instance = WeApp{AppId: c.AppId, AppSecret: c.AppSecret}
	}
	return
}

func (w WeApp) Code2Session(code string) (r Code2SessionResp, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + w.AppId +
		"&secret=" + w.AppSecret +
		"&js_code=" + code +
		"&grant_type=authorization_code"

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}

	if r.ErrorCode != 0 {
		err = errors.New(r.ErrorMsg)
		return
	}

	return
}
