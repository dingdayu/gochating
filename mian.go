package main

import (
	"net/http"
	"text/template"
)

func main() {
	// 注册一个路由
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// 将字符串通过回写指针返回给浏览器
		locals := make(map[string]interface{})
		locals["name"] = "dingdayu"
		t,err := template.ParseFiles("./templates/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, locals)
	})
	// 监听端口 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
