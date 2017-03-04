package handlers

import (
	"net/http"
	"html/template"
)

func Hello(w http.ResponseWriter, r *http.Request)  {
	// 将字符串通过回写指针返回给浏览器
	locals := make(map[string]interface{})
	locals["name"] = "dingdayu"
	t,err := template.ParseFiles("./templates/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, locals)
}
