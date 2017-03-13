package models

import "github.com/dingdayu/gochatting/structs"

func GetUserInfo(token string) *structs.UserInfo  {
	var userInfo *structs.UserInfo = &structs.UserInfo{}
	if token == "1" {
		userInfo = &structs.UserInfo{UID: 1,
			UserName: "dingdayu",
			Email:    "614422099@qq.com"}
	} else if token == "2" {
		userInfo = &structs.UserInfo{UID: 2,
			UserName: "dingxiaoyu",
			Email:    "1003280349@qq.com"}
	}
	return userInfo
}