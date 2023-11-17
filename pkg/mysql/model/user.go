package model

import "strconv"

func (a User) GetUid() string {
	return strconv.Itoa(int(a.UserID))
}
