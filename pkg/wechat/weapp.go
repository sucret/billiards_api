package wechat

import (
	"billiards/pkg/config"
	"billiards/pkg/redis"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

type QrcodeResp struct {
	Buffer    string `json:"buffer"`
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
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

func (w WeApp) GenQrcode(page, scene string) (r QrcodeResp) {
	token, err := w.GetAccessToken()
	if err != nil {
		return
	}

	url := baseUrl + "/wxa/getwxacodeunlimit?access_token=" + token

	param := make(map[string]string)
	param["page"] = page
	param["scene"] = scene
	param["env_version"] = "trial"

	byteParam, _ := json.Marshal(param)
	reader := bytes.NewReader(byteParam)
	resp, err := http.Post(url, "application/json;charset=UTF-8", reader)
	fmt.Println(err)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &r)
	fmt.Println(r)

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

// 微信小程序解密算法 AES-128-CBC
func (w WeApp) DecryptWXOpenData(sessionKey, encryptData, iv string) (map[string]interface{}, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}
	sessionKeyBytes, errKey := base64.StdEncoding.DecodeString(sessionKey)
	if errKey != nil {
		return nil, errKey
	}
	ivBytes, errIv := base64.StdEncoding.DecodeString(iv)
	if errIv != nil {
		return nil, errIv
	}
	dataBytes, errData := aesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)
	fmt.Printf("dataBytes: %v\n", dataBytes)
	if errData != nil {
		return nil, errData
	}

	var result map[string]interface{}
	errResult := json.Unmarshal(dataBytes, &result)

	watermark := result["watermark"].(map[string]interface{})
	if watermark["appid"] != w.AppId {
		return nil, errors.New("Invalid appid data!")
	}
	return result, errResult
}

// AES 解密
func aesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	// 原始数据
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	// 去除填充  --- 数据尾端有'/x0e'占位符,去除它
	length := len(origData)
	unp := int(origData[length-1])
	return origData[:(length - unp)], nil
}
