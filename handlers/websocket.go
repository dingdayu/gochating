package handlers

import (
	"errors"
	"fmt"
	"github.com/dingdayu/gochatting/drives/session"
	"github.com/dingdayu/gochatting/models"
	"github.com/dingdayu/gochatting/structs"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var connectingPool *ConnectingPool = &ConnectingPool{}

func Connection(ws *websocket.Conn) {
	// 检查用户登陆，获取用户对象
	user, err := checkLogin(ws)

	if err != nil {
		needLogin(ws)
		return
	}
	// 检查是否已链接
	if _, ok := connectingPool.Users[user.ID.Hex()]; ok {
		//return
	}

	OnlineUser := &OnlineUser{
		Connection: ws,
		Send:       make(chan structs.Message, 256),
		UserInfo:   &user,
	}

	connectingPool.Users[user.ID.Hex()] = OnlineUser
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
	//uid := onlineUser.UserInfo.ID;
	// target := []structs.UserInfo{
	//	{1,"dingdayu", "614422099@qq.com", "" },
	//	{2,"dingxiaoyu", "1003280349@qq.com", "" },
	//}
	//for _, t := range target{
	//	if t.ID != onlineUser.UserInfo.ID {
	//		m := structs.OnlineNotice(onlineUser.UserInfo.ID, strconv.Itoa(t.ID))
	//		Send(t.ID, m)
	//	}
	//}
	for _, t := range connectingPool.Users {
		if t.UserInfo.ID != onlineUser.UserInfo.ID {
			m := structs.OnlineNotice(onlineUser.UserInfo.ID.Hex(), t.UserInfo.ID.Hex())
			Send(t.UserInfo.ID.Hex(), m)
		}
	}

}

// 有用户退出，将新的用户列表发送给所有人
func (this *OnlineUser) killUserResource() {
	this.Connection.Close()
	id := this.UserInfo.ID
	delete(connectingPool.Users, this.UserInfo.ID.Hex())
	close(this.Send)

	// 用户下线通知，同上面行数逻辑类似
	for _, t := range connectingPool.Users {
		if t.UserInfo.ID != id {
			m := structs.OnlineNotice("0", t.UserInfo.ID.Hex())
			Send(t.UserInfo.ID.Hex(), m)
		}
	}
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
		for _, t := range connectingPool.Users {
			if t.UserInfo.ID != this.UserInfo.ID {
				m := structs.MessageNotice(this.UserInfo.ID.Hex(), t.UserInfo.ID.Hex(), content)
				Send(t.UserInfo.ID.Hex(), m)
			}
		}
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
func Send(uid string, msg *structs.Message) {
	if onlineUser, ok := connectingPool.Users[uid]; ok {
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
	connectingPool = &ConnectingPool{
		Users:     make(map[string]*OnlineUser),
		Broadcast: make(chan structs.Message),
		CloseSign: make(chan bool),
	}
	go connectingPool.run()
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

// 将连接池开放给其他包访问
func GetConnectingPool() *ConnectingPool {
	return connectingPool
}

// 上线用户的连接池
type ConnectingPool struct {
	Users     map[string]*OnlineUser
	Broadcast chan structs.Message
	CloseSign chan bool
}

// 上线用户结构
type OnlineUser struct {
	Connection *websocket.Conn
	Send       chan structs.Message
	UserInfo   *structs.UserInfo
}

// 检查用户登陆
func checkLogin(ws *websocket.Conn) (structs.UserInfo, error) {
	gosessionid, err := ws.Request().Cookie("gosessionid")
	if err != nil {
		needLogin(ws)
		return structs.UserInfo{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}
	sess, _ := session.GlobalSessions.GetSessionStore(gosessionid.Value)
	isLogin := sess.Get("isLogin")
	id := sess.Get("id")
	if isLogin == nil || id == nil || !isLogin.(bool) {
		needLogin(ws)
		return structs.UserInfo{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	user, err := models.IdToUser(id.(string))
	return user, err
}

// 返回需要登陆
func needLogin(ws *websocket.Conn) {
	json := make(map[string]interface{})
	json["code"] = http.StatusUnauthorized
	json["msg"] = http.StatusText(http.StatusUnauthorized)
	websocket.JSON.Send(ws, json)
}
