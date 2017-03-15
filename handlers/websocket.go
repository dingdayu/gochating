package handlers

import (
	"fmt"
	"log"
	"golang.org/x/net/websocket"
	"github.com/dingdayu/gochatting/models"
	"github.com/dingdayu/gochatting/structs"
	"strconv"
)

var ConnectingPool *ConnectingPool = &ConnectingPool{}

func Connection(ws *websocket.Conn) {
	token := ws.Request().URL.Query().Get("token")

	if token == "" {
		return
	}
	user := models.GetUserInfo(token)

	if user == nil {
		return
	}
	// 检查是否已链接
	if _,ok := ConnectingPool.Users[user.UID]; ok {
		//return
	}

	OnlineUser := &OnlineUser{
		Connection: ws,
		Send:       make(chan structs.Message, 256),
		UserInfo:   user,
	}

	ConnectingPool.Users[user.UID] = OnlineUser
	// TODO::Hook
	fmt.Println("新用户上线", user)

	// 通知上线消息：发送消息给固定在线的好友，或聊天室 更新用户在线状态
	go OnlineUser.UserOnlineNotice()

	// 推送队列消息
	go OnlineUser.PushToCline()
	// 等待客户端消息
	OnlineUser.PullFromClient()

	fmt.Println("用户下线", OnlineUser.UserInfo.UserName)
	//用户下线
	OnlineUser.killUserResource()
}

// 通知相应的用户，本用户上线消息
func (onlineUser *OnlineUser) UserOnlineNotice() {
	// 获取所有需要通知到的用户，并通过onlineUser通知过去
	//uid := onlineUser.UserInfo.UID;
	 target := []structs.UserInfo{
		{1,"dingdayu", "614422099@qq.com" },
		{2,"dingxiaoyu", "1003280349@qq.com" },
	}
	for _, t := range target{
		if t.UID != onlineUser.UserInfo.UID {
			m := structs.OnlineNotice(onlineUser.UserInfo.UID, strconv.Itoa(t.UID))
			Send(t.UID, m)
		}
	}


}

// 有用户退出，将新的用户列表发送给所有人
func (this *OnlineUser) killUserResource() {
	this.Connection.Close()
	delete(ConnectingPool.Users, this.UserInfo.UID)
	close(this.Send)

	// 用户下线通知，同上面行数逻辑类似
}

// 等待客户端消息
func (this *OnlineUser) PullFromClient() {
	for {
		var content string
		err := websocket.Message.Receive(this.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		// 如果用户关闭链接，或者刷新浏览器，会出现一个错误
		if err != nil {
			fmt.Println(err)
			return
		}

		//收到客户端消息content
	}
}

// 将Send队列消息推送出去
func (this *OnlineUser) PushToCline() {
	for b := range this.Send {
		err := websocket.JSON.Send(this.Connection, b)
		if err != nil {
			break
		}
	}
}

// 直接发送消息给对于用户
func Send(uid int, msg *structs.Message) {
	if onlineUser, ok := ConnectingPool.Users[uid]; ok {
		err := websocket.JSON.Send(onlineUser.Connection, msg)
		if err != nil {
			log.Println("[ERROR] 消息推送出错！")
		}
	} else {
		// 如果用户不再线
		// 上线后，将消息放到Send队列中
		log.Println("[ERROR] 用户不在线！")
	}

}

func init() {
	ConnectingPool = &ConnectingPool{
		Users:     make(map[int]*OnlineUser),
		Broadcast: make(chan structs.Message),
		CloseSign: make(chan bool),
	}
	go ConnectingPool.run()
}

// 等待通知和关闭通知
func (this *ConnectingPool) run() {
	for {
		select {
		case b := <-this.Broadcast:
			for _, online := range this.Users {
				online.Send <- b
			}
		case c := <-this.CloseSign:
			if c == true {
				close(this.Broadcast)
				close(this.CloseSign)
				return
			}
		}
	}
}

// 上线用户的连接池
type ConnectingPool struct {
	Users     map[int]*OnlineUser
	Broadcast chan structs.Message
	CloseSign chan bool
}

// 上线用户结构
type OnlineUser struct {
	Connection *websocket.Conn
	Send       chan structs.Message
	UserInfo   *structs.UserInfo
}



