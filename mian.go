package main

import (
	"io"
	"net/http"
)

func main() {
	// 注册一个路由
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// 将字符串通过回写指针返回给浏览器
		io.WriteString(w, "hello word!")
	})
	// 监听端口 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
