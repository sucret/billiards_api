package equipment

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 可以读取状态
// 可以修改状态

// http://10.0.69.84/table/1?act=set-status&tm=light
// http://10.0.69.84/table/1?act=status&tm=light
type Terminal struct {
	Url    string
	Status int32
}

type Response struct {
	Table    int
	Terminal string
	Action   string
	Status   int32
}

func (t *Terminal) GetStatus() (res Response, err error) {
	url := t.Url + "&act=status"
	res, err = t.httpGet(url)

	return
}

func (t *Terminal) SetStatus(status int32) (res Response, err error) {
	url := t.Url + "&act=set-status&status=" + strconv.Itoa(int(status))
	res, err = t.httpGet(url)

	return
}

func (t *Terminal) httpGet(url string) (response Response, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	return
}
