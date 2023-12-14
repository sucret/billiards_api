package tool

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Dump(data interface{}) {
	buf, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	fmt.Printf("\n %c[1;0;36m%s%c[0m\n\n", 0x1B, string(buf), 0x1B)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func HttpBuildQuery(params map[string]interface{}, parentKey string) (str string) {
	//fmt.Println("parentKey ", parentKey)
	paramsArr := make([]string, 0)
	for k, v := range params {
		if vals, ok := v.(map[string]interface{}); ok {
			if parentKey != "" {
				k = fmt.Sprintf("%s[%s]", parentKey, k)
			}
			paramsArr = append(paramsArr, HttpBuildQuery(vals, k))
		} else {
			if parentKey != "" {
				paramsArr = append(paramsArr, fmt.Sprintf("%s[%s]=%s", parentKey, k, gconv.String(v)))
			} else {
				paramsArr = append(paramsArr, fmt.Sprintf("%s=%s", k, gconv.String(v)))
			}
		}
	}
	str = strings.Join(paramsArr, "&")
	return
}

func GenerateOrderNum() (orderNum string) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(999999)
	orderNum = time.Now().Format("20060102150405") + fmt.Sprintf("%0*d", 6, num)

	return
}

func Distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const R = 6378.1                // 地球平均半径，单位为千米
	radLat1 := math.Pi * lat1 / 180 // 将角度转换成弧度
	radLon1 := math.Pi * lon1 / 180
	radLat2 := math.Pi * lat2 / 180
	radLon2 := math.Pi * lon2 / 180

	dLat := radLat2 - radLat1
	dLon := radLon2 - radLon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(radLat1)*math.Cos(radLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c

	return distance * 1000 // 返回结果单位为米
}
