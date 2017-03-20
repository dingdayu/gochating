package session

import (
	"github.com/astaxie/beego/session"
	"encoding/json"
)

// 全局session
var GlobalSessions *session.Manager

func init()  {
	// session
	config := `{"cookieName":"gosessionid", "enableSetCookie": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`
	var c session.ManagerConfig
	json.Unmarshal([]byte(config), &c)
	GlobalSessions, _ = session.NewManager("memory", &c)
	go GlobalSessions.GC()
}
