package wechat

import (
	"billiards/pkg/config"
	"billiards/pkg/redis"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const baseUrl = "https://api.weixin.qq.com"

type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	ExpiresIn  int    `json:"expires_in"`
	OpenId     string `json:"openid"`
	ErrorCode  int    `json:"errcode"`
	ErrorMsg   string `json:"errmsg"`
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
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
	url := baseUrl + "/sns/jscode2session?appid=" + w.AppId +
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

// 获取access_token，并缓存7000秒
// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func (w WeApp) GetAccessToken() (token string, err error) {
	key := "access_token:" + w.AppId

	token = redis.GetRedis().Get(key).Val()

	if token == "" {
		url := baseUrl + "/cgi-bin/token?grant_type=client_credential&appid=" + w.AppId + "&secret=" + w.AppSecret

		resp, err := http.Get(url)
		if err != nil {
			err = errors.New("获取access_token失败")
			return token, err
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		body, _ := ioutil.ReadAll(resp.Body)

		res := AccessTokenResp{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			return token, err
		}

		redis.GetRedis().Set("access_token:"+w.AppId, res.AccessToken, time.Second*7000)

		token = res.AccessToken
	}

	return
}
