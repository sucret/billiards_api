package tool

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"os"
	"strings"
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
