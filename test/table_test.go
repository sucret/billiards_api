package test

import (
	"billiards/pkg/wechat"
	"testing"
)

func TestGenTableQr(t *testing.T) {
	wechat.GetApp("client").GenQrcode("pages/shop/list/shopList", "a=1")
}
