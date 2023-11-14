package response

import "gin-api/pkg/mysql/model"

type AdminDetail struct {
	model.Admin
	RoleInfo []struct {
		RoleId uint   `json:"role_id"`
		Name   string `json:"name"`
	} `json:"role_info"`
	RoleList []model.Role `json:"role_list"`
	NodeList []NodeTree   `json:"node_list"`
}

type NodeMeta struct {
	Icon         string `json:"icon"`
	Show         bool   `json:"show"`
	Title        string `json:"title"`
	HideChildren bool   `json:"hideChildren"`
}

type AdminNode struct {
	Component string   `json:"component"`
	Id        int32    `json:"id"`
	Meta      NodeMeta `json:"meta"`
	Name      string   `json:"name"`
	ParentId  int32    `json:"parentId"`
	Redirect  string   `json:"redirect"`
	Path      string   `json:"path"`
}
