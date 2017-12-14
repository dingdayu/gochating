package api

import (
	"net/http"
	"github.com/dingdayu/gochatting/drives/session"
	"github.com/dingdayu/gochatting/utils"
	"github.com/dingdayu/gochatting/models"
)

func Register(w http.ResponseWriter, r *http.Request)  {
	sess, _ := session.GlobalSessions.SessionStart(w, r)

	ret := make(map[string]interface{})
	// 这里不遵循大写字母开头的问题

	username := r.FormValue("username")
	passwd := r.FormValue("passwd");
	email := r.FormValue("email");

	if username == "" {
		ret["code"] = 301
		ret["msg"] = "username not empty!"
		utils.ReturnJson(ret, w)
		return
	}
	if passwd == "" {
		ret["code"] = 302
		ret["msg"] = "passwd not empty!"
		utils.ReturnJson(ret, w)
		return
	}

	visitLogin := sess.Get("visitLogin")
	if visitLogin == nil || !visitLogin.(bool) {
		ret["code"] = 304
		ret["msg"] = "login error!"
		utils.ReturnJson(ret, w)
		return
	}

	userInfo, err := models.UsernameToUser(username)
	if err != nil{
		ret["code"] = 305
		ret["msg"] = "username error!"
		utils.ReturnJson(ret, w)
		return
	}

	models.AddUser(username,passwd, email)

	sess.Set("isLogin", true)
	sess.Set("id", userInfo.ID.Hex())
	ret["code"] = 200
	ret["msg"] = "success"
	utils.ReturnJson(ret, w)
	sess.Set("isLogin", true)
	return

}