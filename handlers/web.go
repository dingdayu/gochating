package handlers

import (
	"net/http"
	"html/template"
	"os"
)

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	// 直接调用http包提供的文件服务方法，直接根据请求路径返回文件内容
	if exists := os.IsExist(r.URL.Path[len("/"):]); !exists {
		// 文件找不到返回404
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, r.URL.Path[len("/"):])
}

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
