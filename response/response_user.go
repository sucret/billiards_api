package response

import (
	"billiards/pkg/mysql/model"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type UserLogin struct {
	User      model.User  `json:"user"`
	TokenData TokenOutPut `json:"token_data"`
}
