package structs

import "time"

const (
	MSG_VERSION = 1
)

// 消息结构
type Message struct {
	Module string `json:"module"`
	// 协议版本号
	Version int `json:"varsion"`
	// 消息唯一id
	UUID int
	// 消息来源用户id
	UID string
	// 消息类型 GROUP PRIVATE SYSTEM
	Type string `json:"type"`
	// 目标id ： ueser_id || group_id
	Target string `json:"target"`
	// 消息发送时间
	Time    int `json:"time"`
	Content MessageContent
}

// 消息内容结构
type MessageContent struct {
	Type string
	Data string
}

func MessageNotice(uid string, target string, msg string) *Message {
	m := Message{
		Module:  "PUSH",
		Version: MSG_VERSION,
		UUID:    00,
		UID:     uid,
		Type:    "SYSTEM",
		Target:  target,
		Time:    int(time.Now().Unix()),
		Content: MessageContent{
			Type: "platform",
			Data: msg,
		},
	}
	return &m
}

// 上线通知
func OnlineNotice(uid string, target string, user *UserInfo) *Message {
	m := Message{
		Module:  "PUSH",
		Version: MSG_VERSION,
		UUID:    00,
		UID:     uid,
		Type:    "SYSTEM",
		Target:  target,
		Time:    int(time.Now().Unix()),
		Content: MessageContent{
			Type: "OnlineNotice",
			Data: user.UserName,
		},
	}
	return &m
}

// 下线通知
func OfflineNotice(uid string, target string, user *UserInfo) *Message {
	m := Message{
		Module:  "PUSH",
		Version: MSG_VERSION,
		UUID:    00,
		UID:     uid,
		Type:    "SYSTEM",
		Target:  target,
		Time:    int(time.Now().Unix()),
		Content: MessageContent{
			Type: "OfflineNotice",
			Data: user.UserName,
		},
	}
	return &m
}
