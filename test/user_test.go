package test

import (
	"billiards/pkg/wechat"
	"fmt"
	"testing"
)

// 发放优惠券
func TestGetUserInfo(t *testing.T) {
	sessionKey := "986ZLRRrjTbtH2GCCM6cig=="
	//e := "61qC0KW+P7CL8kl3Cw3ArY+sqtfw+ks5lUx0I29lRAAwWJgpa1ZEV1FM1AV22iHtKPbsRhjdIko73nRXsmuLkt0G2QFrVb33md0A6lzGGRK9NT8Vf9nUM2ylB43Oa2ySGbDM4cee7WM7mV97V/H9prxr7U0oCRO9H2IaMZIlsosMVslWi4Tka1vwk4OdRcqGoVxMyhnXRm2CRkkGvO7EreB7x+cjKP6cVmJVkaWI/YthpgniGfjszBwvvgUAgtj4PK0NuQo0mfShU7Yg3tBCGCQux3ptFabfCTtWhMBRpWhSE0/X3s3PDHOHfArid8udksyrBFzOO7MgcuhGHa/AOAwkjnHCUi+HNYCQzpLp7s4XquC7wxFk7mFntkbjoTiBLEcexF+jNppQX8WLuXn/iQdMIp1Rcq6uSNry33Ezg0RkWN+3+2zMgAF0xB/iIdzIDmOiaocTSDtAHRdnHIdwcA=="

	e := "9a0HnugHzr0qfGNABL2s3athdzBFV+EySv4q5+ZmFLfQTfvgWLYNp/XHkyOKWHPO9r08aqDhNFLQcqyPeqRwW9UP7PsCEecR4bSsY4TsSPh8HJilcFEBBvS5kBWgBET4y8Xewx8OlGpZGkn4kP/4M9HNG5jjlWl3N8v+oGuNIMYkgMh7RFid7QqJCawEjC904sXu3HJiSDrqlWARFRx41Tdue6cFdDUVbVBkBw11eHMKmihqFdE+eKh5/4/PgtxAHqILCPUB57bFWoaAMQLFJJamNg+EN+WA/jgAOzUIC760ti8P96qLI0VLpQrWieyV2C+ltYfU/c32BF4VZjuLongqfWw540kGc1TMYYzF5w/li07ZsqmeLo2FJgs7yRfdVPV8tQgrm1BGobtRE6rWozwp/iqn5LsMjS8vrtPdx6ldpLUlep9mbWQIOcajaLJooY9PwrlQ9H8jnOEWxqNZbQ=="

	w := wechat.GetApp("client")

	iv := "9yabobBxun/ySq+mLUF/Og=="
	m, err := w.DecryptWXOpenData(sessionKey, e, iv)

	fmt.Println(m, err)

	fmt.Printf("解密后的数据: %v\n", m)
	//uc, err := service.UserCouponService.HandOut(3, 4)
	//fmt.Println(uc, err)
}
