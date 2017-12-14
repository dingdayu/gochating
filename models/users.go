package models

import (
	"fmt"
	"github.com/dingdayu/gochatting/drives/db"
	"github.com/dingdayu/gochatting/structs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DB_NAME  = "chat"

// 通过id查询用户
func IdToUser(id string) (result structs.UserInfo, err error) {

	mgoSess := db.CloneSession()//调用这个获得session
	defer mgoSess.Close()  //一定要记得释放

	c := mgoSess.DB(DB_NAME).C("user")
	result = structs.UserInfo{}
	where := bson.M{"_id": bson.ObjectIdHex(id)}
	err = c.Find(where).One(&result)
	if err != nil {
		return
	}
	return
}

// 通过username查询用户
func UsernameToUser(username string) (result structs.UserInfo, err error) {

	mgoSess := db.CloneSession()//调用这个获得session
	defer mgoSess.Close()  //一定要记得释放

	c := mgoSess.DB(DB_NAME).C("user")
	result = structs.UserInfo{}
	where := bson.M{"username": username}
	err = c.Find(where).One(&result)
	if err != nil {
		return
	}
	return
}

func AddUser(u,p,e string)  {
	mgoSess := db.CloneSession()//调用这个获得session
	defer mgoSess.Close()  //一定要记得释放

	c := mgoSess.DB(DB_NAME).C("user")
	err := c.Insert(structs.UserInfo{
		ID:bson.NewObjectId(),
		UserName:u,
		Passwd:p,
		Email:e,
	})

	if err != nil {
		panic(err)
	}
}

func GetUser() {
	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB_NAME).C("user")
	c.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	//err = c.Insert(&structs.UserInfo{1,"dingdayu", "614422099@qq.com",bson.NewObjectId()})
	//if err != nil { panic(err) }
	result := structs.UserInfo{}
	err = c.Find(bson.M{"username": "dingdayu"}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phone:", result)
}

func VerifyLogin(username string, passwd string) bool {
	c := db.GetSession().DB(DB_NAME).C("user")
	result := structs.UserInfo{}
	where := bson.M{"username": username, "passwd": passwd}
	err := c.Find(where).One(&result)
	if err != nil {
		return false
	}
	if username == result.UserName {
		return true
	} else {
		return false
	}
}
