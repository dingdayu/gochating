package api

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io"
	"github.com/dingdayu/gochatting/structs"
	"github.com/dingdayu/gochatting/handlers"
)

func HelloJson(w http.ResponseWriter, r *http.Request) {
	// 定义返回的结构体
	type jsonType struct {
		// 这里遵循大写字母开头方可被开放
		// 原因在于自定义结构体里面的对象，需要可以被json包访问到，
		// 而go规定只有大写开头的才能被包外部访问，而类型属于go语言的基本结构
		Name string
		age  int
	}

	// 实例化一个结构体
	hello := jsonType{Name: "dingdayu", age: 23}
	// map类型同样的使用方法
	//hello := make(map[string]string)
	// 这里不遵循大写字母开头的问题
	//hello["Name"] = "dingdayu"
	//hello["age"] = 23

	// 将结构体或类型转json字符串 除channel,complex和函数几种类型外，都可以转json
	// 注意  json.Marshal() 返回的是字节 需要转 string()
	if j, err := json.Marshal(hello); err != nil {
		fmt.Fprint(w, "json error")
	} else {
		// 返回json的类型头信息
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(j))
	}
}

func GetOnlineUserList(w http.ResponseWriter, r *http.Request)  {
	// 先获取所有在线的人
	var userList []structs.UserInfo

	ConnectingPool := handlers.GetConnectingPool()
	for _, online := range ConnectingPool.Users {
		userList = append(userList, *online.UserInfo)
	}

	if j, err := json.Marshal(userList);  userList == nil || err != nil {
		fmt.Fprint(w, "json error")
		io.WriteString(w, "[]")
	} else {
		// 返回json的类型头信息
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(j))
	}
}