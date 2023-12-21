package router

import (
	api "billiards/api/external"
	"billiards/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func setExternalRoute(r *gin.Engine) {
	clientRouters := r.Group("/e").Use(gin.Logger(), middleware.CustomRecovery())
	{
		// 微信支付回调
		clientRouters.POST("/wechat/pay-notify", api.WechatApi.PayNotify)

		// 退款测试
		//clientRouters.GET("/wechat/refund", api.WechatApi.Refund)

		clientRouters.GET("/shop/terminal/status", api.ShopApi.TerminalStatus)

		// 店铺终端状态socket
		clientRouters.GET("/shop/terminal-status/socket", api.ShopApi.TerminalStatusSocket)
		clientRouters.GET("/shop/chan-test", api.ShopApi.ChanTest)

		clientRouters.GET("/shop/terminal/state", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer func(conn *websocket.Conn) {
				_ = conn.Close()
			}(conn)

			//time.AfterFunc(4*time.Second, func() {
			//	err = conn.WriteMessage(websocket.TextMessage, []byte("已建立连接"))
			//	if err != nil {
			//		return
			//	}
			//})

			// 建立一个映射店铺id的chan类型的map
			// 如果店铺有变动就往chan中推送店铺的信息
			// 在这里监控chan

			for {
				_, msg, err := conn.ReadMessage()

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Printf("收到消息:%s \n", msg)

				err = conn.WriteMessage(websocket.TextMessage, []byte("已收到消息"))
				if err != nil {
					fmt.Println(err)
				}
			}
		})
	}
}
