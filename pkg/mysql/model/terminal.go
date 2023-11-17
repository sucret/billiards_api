package model

const (
	TerminalStatusOpen  = 1 + iota // 开启
	TerminalStatusClose            // 关闭
)

const (
	TerminalTypeLight       = 1 + iota // 射灯
	TerminalTypeLocker                 // 锁球器
	TerminalTypePicReader              // 图像识别芯片
	TerminalTypeShopCamera             // 店铺摄像头
	TerminalTypeTableCamera            // 球桌摄像头
	TerminalTypeShopLight              // 店铺灯光
)
