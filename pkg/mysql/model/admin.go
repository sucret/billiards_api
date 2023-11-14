package model

import "strconv"

func (a Admin) GetUid() string {
	return strconv.Itoa(int(a.AdminID))
}
