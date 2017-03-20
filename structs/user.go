package structs

import "gopkg.in/mgo.v2/bson"


// 用户的资料信息结构
type UserInfo struct {
	UID      int
	UserName string `bson:"username"`
	Email    string `bson:"email"`
	ID bson.ObjectId `bson:"_id"`
	Passwd string `bson:"passwd"`
}

